$session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$session.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"
$session.Cookies.Add((New-Object System.Net.Cookie(".AspNetCore.Culture", "c%3Dzh-Hans%7Cuic%3Dzh-Hans", "/", ".dangquyun.com")))
Invoke-WebRequest -UseBasicParsing -Uri "https://app.dangquyun.com/apis/tenancy/api/app/applicationTemplate/generateAppTemplate" `
-Method "POST" `
-WebSession $session `
-Headers @{
"authority"="app.dangquyun.com"
  "method"="POST"
  "path"="/apis/tenancy/api/app/applicationTemplate/generateAppTemplate"
  "scheme"="https"
  "__tenant"="3a01177d-19e6-b5bc-56d6-74fcca16f4dc"
  "accept"="application/json"
  "accept-encoding"="gzip, deflate, br"
  "accept-language"="en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7"
  "authorization"="Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IkRDRTlEQjM2RjZBRkEwNUUyMzlEOUFDNjFCQTcyQ0M4QzEzNDJCNjEiLCJ0eXAiOiJKV1QiLCJ4NXQiOiIzT25iTnZhdm9GNGpuWnJHRzZjc3lNRTBLMkUifQ.eyJzY29wZSI6WyJvcGVuaWQiLCJ0ZW5hbnRfdXNlciJdLCJfX3RlbmFudCI6IjNhMDExNzdkLTE5ZTYtYjViYy01NmQ2LTc0ZmNjYTE2ZjRkYyIsInN1YiI6ImQ0Yjk5ZmFkLTI5MmItNDNkNi04OTAwLWI1MDY4MDQxOGViOSIsImF1ZCI6InB0X2FwaXMiLCJleHAiOjE2NjAyODgyODQsImlzcyI6Imh0dHBzOi8vYWNjb3VudC5kYW5ncXV5dW4uY29tL3VzZXJpZGVudGl0eSIsImlhdCI6MTY2MDI2NjY4NCwibmJmIjoxNjYwMjY2Njg0fQ.LQrp8DBg611nB87Pj95Yps-KvkZe0gz1HFhjZwbQLUAy6fwkKlx9jkkMD-bAvWrNafF63ucgLkXTdkZXzDjsPDaxIm_mVlb6KR_6mJgF8YT2Mj0bLMRQNZmBmhnaDMi9_WXkfhSPV_olZzy6nYcQxmEnqEYAElsseZktOMJbFJM8XDr1u57omlf3NknnoJoFd5RI3UgMwy8vvwRIToxWlo_huMW9baSXUgVib5aLH8XlXjZxsX12xiIp4M8QmyR0YuAoHcu0RSHctH9qoDdVHHH5W6xOGn3D85izcRKX2ulp6rKpbmA1VObIMkceGJiDvmZQpWSYIKKrJF_KYM_e6Q"
  "cache-control"="no-cache"
  "origin"="https://app.dangquyun.com"
  "pragma"="no-cache"
  "referer"="https://app.dangquyun.com/tabs/39fed0d0-21bb-97e2-b814-3f92669cc43c/dynamicApp/3a011873-4795-192f-8be9-b488d8a23f36/pagedetail/3a002480-dfcb-2c88-ed1b-52ad3c00d3ac/host"
  "sec-ch-ua"="`"Chromium`";v=`"104`", `" Not A;Brand`";v=`"99`", `"Google Chrome`";v=`"104`""
  "sec-ch-ua-mobile"="?0"
  "sec-ch-ua-platform"="`"Windows`""
  "sec-fetch-dest"="empty"
  "sec-fetch-mode"="cors"
  "sec-fetch-site"="same-origin"
} `
-ContentType "application/json;charset=UTF-8" `
-Body "{`"appId`":`"3a011873-4795-192f-8be9-b488d8a23f36`",`"roles`":[]}"