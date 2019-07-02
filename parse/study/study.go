package study

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync_study/configuration"
	"sync_study/entity"
	"sync_study/logger"
	"sync_study/utils"
)

var (
	regexLink = regexp.MustCompile(`<a href="(.*[^"])">(.*[^<])</a>`)
)

// ParseStudy -
func ParseStudy(result *entity.RequestResult) *entity.ParseRequest {

	var (
		requests     entity.ParseRequest
		studyRequest *entity.StudyRequest
		body         []byte
		resp         *http.Response
		ok           bool
		err          error
	)

	if result == nil {
		return nil
	}

	source, err := url.Parse(result.SourceURI)
	if err != nil {
		logger.Errorf("sourceURI[%s]错误[%s]", result.SourceURI, source)
		return nil
	}

	if resp, ok = result.Body.(*http.Response); ok {
		defer resp.Body.Close()
		if resp.ContentLength > 1<<10*10<<10 {
			// 保存文件 大于10M
			filename := configuration.Config().DataPath["filePath"] + source.Path
			file, err := utils.OpenFile(filename)
			if err != nil {
				logger.Errorf("打开文件[%s]失败[%s]", filename, err)
				return nil
			}

			logger.Infof("下载内容到文件[%s]", filename)
			written, err := io.Copy(file, resp.Body)
			if err != nil {
				file.Close()
				os.Remove(filename)
				logger.Errorf("写入内容到文件[%s]失败[%s]", filename, err)
				return nil
			}

			if written < resp.ContentLength {
				file.Close()
				os.Remove(filename)
				logger.Errorf("文件[%s]写入不完整", filename)
				return nil
			}

			defer file.Close()
			return nil
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Errorf("读取接口[%s]返回的结果失败~%s", source.String(), err)
			return nil
		}
	} else {
		if body, ok = result.Body.([]byte); !ok || len(body) < 1 {
			logger.Errorf("URL[%s]未知的返回结果~", source.String())
			return nil
		}
	}

	// logger.Infof("body: %s", body)

	if err == nil {
		if !utils.CheckURLPath(source) {
			// 异步写入到文件
			go func() {
				filename := configuration.Config().DataPath["filePath"] + source.Path
				utils.WriterFile(filename, body, true)
				logger.Infof("保存文件[%s]", filename)
			}()
			return nil
		}
	}

	links := regexLink.FindAllStringSubmatch(string(body), -1)

	if links == nil || len(links) < 1 {
		return nil
	}

	if sr, ok := result.Data.(*entity.StudyRequest); ok {
		studyRequest = sr
		studyRequest.Depth++
	}

	for i := range links {
		if len(links[i]) != 3 {
			continue
		}
		pathName := links[i][2]
		if pathName == "../" {
			continue
		}
		uri := result.SourceURI + links[i][1]

		// 是否已存在   重复代码  用于加快重启后的数据同步
		if uriSource, err := url.Parse(uri); err == nil {
			if !utils.CheckURLPath(uriSource) { // 判断是否为目录
				filename := configuration.Config().DataPath["filePath"] + uriSource.Path
				if utils.FileExists(filename) {
					// 文件已存在
					logger.Infof("file: %s 已存在", filename)
					continue
				}
			}
		}

		requests.Items = append(requests.Items,
			&entity.Request{
				URI:           uri,
				ParseHandler:  ParseStudy,
				RequestConfig: result.RequestConfig,
				Data:          studyRequest,
			},
		)
	}

	return &requests
}
