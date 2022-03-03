package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Job struct {
	JobName     string `json:"job_name"`
	JobCity     string `json:"job_city"`
	JobCategory string `json:"job_category"`
	JobType     string `json:"job_type"`
	JobDescribe string `json:"job_describe"`
	JobDemand   string `json:"job_demand"`
}

var param = map[string]int{
	"Limit":  0,
	"Offset": 0,
}

var result map[string]interface{}
var err error

func main() {
	// 分两次调用接口 第一次求count职位数 第二次求所有职位信息

	// 第一次 将count赋值给Limit
	if result, err = SendPostRequest(); err != nil {
		log.Println("error:", err)
		return
	}
	r := result["data"].(map[string]interface{})
	param["Limit"] = int(r["count"].(float64))

	// 第二次 求所有职位信息
	if result, err = SendPostRequest(); err != nil {
		log.Println("error:", err)
		return
	}
	data := result["data"].(map[string]interface{})
	job_post_list := data["job_post_list"].([]interface{})
	jobList := make([]Job, 0)
	for index, _ := range job_post_list {
		jobMap := job_post_list[index].(map[string]interface{})
		job := Job{
			JobName:   jobMap["title"].(string),
			JobCity:   jobMap["city_info"].(map[string]interface{})["name"].(string),
			JobDemand: jobMap["requirement"].(string),
			JobType:   jobMap["job_category"].(map[string]interface{})["name"].(string),
			//JobCategory: 	jobMap["job_category"].(map[string]interface{})["parent"].(map[string]interface{})["name"].(string),
			JobDescribe: jobMap["description"].(string),
		}
		jobList = append(jobList, job)
	}

	//生成写入json格式文件
	b, err := json.Marshal(jobList)
	//log.Println(string(b))
	if err != nil {
		log.Println("error:", err)
		return
	}
	if err = ioutil.WriteFile("job.json", b, 0666); err != nil {
		log.Println("error:", err)
		return
	}

	log.Println("生成json文件成功！")
}

func SendPostRequest() (result map[string]interface{}, err error) {
	// 后端返回列表接口地址
	url := "https://xskydata.jobs.feishu.cn/api/v1/search/job/posts"

	client := &http.Client{}

	// 序列化请求参数
	p, err := json.Marshal(param)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(p))

	// 分别设置三个请求头
	// 第一个请求头为固定是用于显示校招内容
	request.Header.Add("website-path", "school")
	// 第二第三个请求头会随token请求变动而变动 其中x-csrf-token为cookie的内置参数 需要具体找到cookie生成逻辑
	FindToken(request)

	response, _ := client.Do(request)
	body, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return result, nil
}

func FindToken(request *http.Request) {
	// 根据浏览器清除缓存后可以发现第一次调用时 返回招聘列表接口会调用两次 其中第一次失败第二次成功
	// 需要辨别cookie生成是通过js方式还是接口方式
	// 于是可以找到处理cookie逻辑接口 https://xskydata.jobs.feishu.cn/api/v1/csrf/token
	tokenUrl := "https://xskydata.jobs.feishu.cn/api/v1/csrf/token"
	req, _ := http.NewRequest("POST", tokenUrl, nil)
	client := &http.Client{}
	rep, _ := client.Do(req)

	csrf_token := strings.Split(rep.Header.Get("set-cookie"), ";")[0]
	cookie := "atsx-portal-session-v1=; channel=saas-career; platform=pc; s_v_web_id=verify_l07iala9_XuJ4GZpg_Vkpl_4mvh_8QdX_Qe9QWFD9KSVM; device-id=7069951796982941195; tea_uid=7070587498893002255; atsx-portal-session-v1=; atsx-portal-user-source-v1=; SLARDAR_WEB_ID=8474f711-475c-4863-8812-610e5729e46e;" + csrf_token
	// 注意csrf_token最后三位 %3D 需修改编码为等于号 =
	csrf_token = strings.Replace(csrf_token, "%3D", "=", -1)
	csrf_token = strings.Split(csrf_token, "atsx-csrf-token=")[1]
	//log.Println(csrf_token,cookie)
	request.Header.Add("x-csrf-token", csrf_token)
	request.Header.Add("cookie", cookie)
}
