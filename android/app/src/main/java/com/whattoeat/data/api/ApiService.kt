package com.whattoeat.data.api

import com.whattoeat.data.api.models.*
import retrofit2.Response
import retrofit2.http.*

interface ApiService {

    // 健康检查
    @GET("health")
    suspend fun health(): Response<HealthResponse>

    // 认证
    @POST("api/auth/register")
    suspend fun register(@Body request: RegisterRequest): Response<ApiResponse<User>>

    @POST("api/auth/login")
    suspend fun login(@Body request: LoginRequest): Response<ApiResponse<LoginResponse>>

    // 游客登录
    @POST("api/auth/guest")
    suspend fun guestLogin(): Response<ApiResponse<LoginResponse>>

    // 菜单
    @GET("api/menus")
    suspend fun getMenus(): Response<ApiResponse<List<Menu>>>

    @POST("api/menus")
    suspend fun createMenu(@Body request: CreateMenuRequest): Response<ApiResponse<CreateMenuResponse>>

    @DELETE("api/menus/{id}")
    suspend fun deleteMenu(@Path("id") id: Long): Response<ApiResponse<Unit>>

    // 餐厅
    @GET("api/restaurants")
    suspend fun getRestaurants(): Response<ApiResponse<List<Restaurant>>>

    // 决策
    @POST("api/decide")
    suspend fun decide(@Body request: DecideRequest): Response<ApiResponse<DecideResponse>>

    @GET("api/history")
    suspend fun getHistory(): Response<ApiResponse<HistoryResponse>>
}
