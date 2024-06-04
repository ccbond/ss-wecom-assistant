// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package server

// setupRouters sets up routers.
func setupAPIRouters(srv *Server) {
	routerGroup := srv.apiServer.GetAPIRouteGroup()
	api := routerGroup.Group("/api/v1")

	wechatApi := api.Group("/wx")
	{
		wechatApi.GET("/auto_reply", srv.wechatCheck)
		wechatApi.POST("/auto_reply", srv.wechatReply)
		wechatApi.GET("/history", srv.getHistory)
		wechatApi.GET("/history_json", srv.getHistoryJson)
	}

	manageApi := api.Group("/manage")
	{
		manageApi.GET("/access_token", srv.getAccessToken)
	}
}
