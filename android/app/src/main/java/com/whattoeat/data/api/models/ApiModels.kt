package com.whattoeat.data.api.models

import com.google.gson.annotations.SerializedName

// 通用响应
data class ApiResponse<T>(
    val code: Int,
    val message: String,
    val data: T?
)

// 用户
data class User(
    val id: Long,
    val username: String
)

// 餐厅
data class Restaurant(
    val id: Long,
    val name: String,
    @SerializedName("created_at") val createdAt: String,
    @SerializedName("updated_at") val updatedAt: String
)

// 菜单
data class Menu(
    val id: Long,
    @SerializedName("restaurant_id") val restaurantId: Long,
    @SerializedName("dish_name") val dishName: String,
    @SerializedName("created_at") val createdAt: String,
    @SerializedName("updated_at") val updatedAt: String,
    val restaurant: Restaurant?
)

// 决策记录
data class DecisionRecord(
    val id: Long,
    @SerializedName("user_id") val userId: Long,
    @SerializedName("menu_id") val menuId: Long,
    @SerializedName("decided_at") val decidedAt: String,
    val menu: Menu?
)

// 登录响应
data class LoginResponse(
    val token: String,
    @SerializedName("user_id") val userId: Long,
    val username: String
)

// 决策响应
data class DecideResponse(
    val menu: Menu,
    val message: String
)

// 历史响应
data class HistoryResponse(
    val records: List<DecisionRecord>,
    val total: Int
)

// 创建菜单响应
data class CreateMenuResponse(
    val menu: Menu,
    @SerializedName("is_new_restaurant") val isNewRestaurant: Boolean
)

// 请求体
data class LoginRequest(
    val username: String,
    val password: String
)

data class RegisterRequest(
    val username: String,
    val password: String
)

data class CreateMenuRequest(
    @SerializedName("restaurant_name") val restaurantName: String,
    @SerializedName("dish_name") val dishName: String
)

data class DecideRequest(
    @SerializedName("menu_ids") val menuIds: List<Long> = emptyList()
)
