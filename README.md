# WFiber
Wfiber stands for **W**rapped **Fiber** is a wrapper over the web framework [gofiber](https://github.com/gofiber/fiber) that has useful code generation and other features such as

## 🎯 Features overview
- Fully typesafe api client generation for the frontend in typescript
    - for endpoints we can define a input and ouput struct and typescript interfaces for the same are generated with the help of [gos2tsi](https://github.com/N4r35h/gos2tsi) which is based on https://golang.org/x/tools/go/packages and doesnt use reflections hence is faster that traditional libraries that generate typescript interfaces from a given golang struct

## 👀 Examples

## ⚡️ Quickstart
install go module
```
go get github.com/N4r35h/wfiber
```