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
	}

	manageApi := api.Group("/manage")
	{
		manageApi.GET("/access_token", srv.getAccessToken)
		manageApi.GET("/msg_history", srv.getMsgHistory)
	}
}
