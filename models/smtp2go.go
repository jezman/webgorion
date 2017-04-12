package models

//func Report() (rerpBody *smtp2go.EmailSummary) {
//    values := map[string]string{"api_key": Conf.Smtp2goApiKey}
//    jsonValue, _ := json.Marshal(values)
//    res, err := http.Post("https://api.smtp2go.com/v3/stats/email_summary", "application/json", bytes.NewBuffer(jsonValue))
//    if err != nil {
//        fmt.Println(err.Error())
//    }
//    body, err := ioutil.ReadAll(res.Body)
//    if err != nil {
//        fmt.Println(err.Error())
//    }
//    var respBody = new(smtp2go.EmailSummary)
//    err = json.Unmarshal(body, &respBody)
//    if err != nil {
//        fmt.Println("Error:", err)
//    }
//    return respBody
//}
