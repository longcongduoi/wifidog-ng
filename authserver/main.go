/*
 * Copyright (C) 2017 Jianhui Zhao <jianhuizhao329@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
    "flag"
    "log"
    "fmt"
    "time"
    "math/rand"
    "strconv"
    "net/http"
    "crypto/md5"
    "encoding/hex"
)

var loginPage = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>WiFi Portal认证登录</title>
    <meta name="viewport" content="width=device-width,minimum-scale=1.0,maximum-scale=1.0,user-scalable=no" />
    <script>
        function load() {
            document.forms[0].action = "/wifidog/login" + window.location.search;
        }
    </script>
    <style type="text/css">
        html, body {
            background-color: #555;
        }

        html {
            height: 100%;
        }

        body {
            height: 98%;
        }

        #box {
            position: absolute;
            top: 50%;
            left:50%;
            margin: -150px 0 0 -150px;
            width: 300px;
            height: 300px;
        }

        #box h1 {
            color: #fff;
            text-shadow:0 0 10px;
            letter-spacing: 1px;
            text-align: center;
        }

        #box input {
            width: 278px;
            height: 18px;
            margin-bottom: 10px;
            outline: none;
            padding: 10px;
            font-size: 15px;
            color: #fff;
            text-shadow:1px 1px 1px;
            border-top: 1px solid #312E3D;
            border-left: 1px solid #312E3D;
            border-right: 1px solid #312E3D;
            border-bottom: 1px solid #56536A;
            border-radius: 4px;
            background-color: #2D2D3F;
        }

        .btn {
            width: 300px;
            min-height: 20px;
            background-color: #4a77d4;
            border: 1px solid #3762bc;
            color: #fff;
            padding: 9px 14px;
            font-size: 20px;
            border-radius: 5px;
        }
    </style>
</head>
<body onload="load()">
    <div id="box">
        <h1>Login</h1>
        <form method="POST">
            <input type="text" required="required" placeholder="username" name="username"></input>
            <input type="password" required="required" placeholder="password" name="password"></input>
            <button class="btn" type="submit">Login</button>
        </form>
    </div>
</body>
</html>
`
var portalPage = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>WiFi Portal</title>
    <meta name="viewport" content="width=device-width,minimum-scale=1.0,maximum-scale=1.0,user-scalable=no" />
    <style type="text/css">
            html, body {
            background-color: #555;
        }

        html {
            height: 100%;
        }

        body {
            height: 98%;
        }

        #box {
            position: absolute;
            top: 50%;
            left:50%;
            margin: -150px 0 0 -150px;
            width: 300px;
            height: 300px;
        }

        #box h1 {
            color: #fff;
            text-shadow:0 0 10px;
            letter-spacing: 1px;
            text-align: center;
        }

        #box input {
            width: 278px;
            height: 18px;
            margin-bottom: 10px;
            outline: none;
            padding: 10px;
            font-size: 15px;
            color: #fff;
            text-shadow:1px 1px 1px;
            border-top: 1px solid #312E3D;
            border-left: 1px solid #312E3D;
            border-right: 1px solid #312E3D;
            border-bottom: 1px solid #56536A;
            border-radius: 4px;
            background-color: #2D2D3F;
        }

        .btn {
            width: 300px;
            min-height: 20px;
            background-color: #4a77d4;
            border: 1px solid #3762bc;
            color: #fff;
            padding: 9px 14px;
            font-size: 20px;
            border-radius: 5px;
        }
    </style>
</head>
<body>
    <div id="box">
        <h1>Welcome to WiFi Portal</h1>
    </div>
</body>
</html>
`
func generateToken(mac string) string {
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(mac + strconv.FormatFloat(rand.Float64(), 'e', 6, 32)))
    cipherStr := md5Ctx.Sum(nil)
    return hex.EncodeToString(cipherStr)
}

func main() {
    port := flag.Int("port", 8912, "http service port")

    flag.Parse()

    rand.Seed(time.Now().Unix())

    http.HandleFunc("/wifidog/ping", func(w http.ResponseWriter, r *http.Request) {
        log.Println("ping", r.URL.RawQuery)
        fmt.Fprintf(w, "Pong")
    })

    http.HandleFunc("/wifidog/login", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            fmt.Fprintf(w, loginPage)
        } else {
            gw_address := r.URL.Query().Get("gw_address")
            gw_port := r.URL.Query().Get("gw_port")
            mac := r.URL.Query().Get("mac")
            token := generateToken(mac)
        
            uri := fmt.Sprintf("http://%s:%s/wifidog/auth?token=%s", gw_address, gw_port, token)
            fmt.Println("Redirect:", uri)
            http.Redirect(w, r, uri, http.StatusFound)
        }
    })

    http.HandleFunc("/wifidog/auth", func(w http.ResponseWriter, r *http.Request) {
        log.Println("auth", r.URL.RawQuery)
        fmt.Fprintf(w, "Auth: 1")
    })

    http.HandleFunc("/wifidog/portal", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, portalPage)
    })

    log.Println("Listen on: ", *port, "SSL off")
    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(*port), nil))
}