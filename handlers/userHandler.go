package handlers

import (
	"audioPhile/database"
	"audioPhile/database/helper"
	"audioPhile/models"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/http"
)

func AddProductToCart(writer http.ResponseWriter, request *http.Request) {
	var cartItem models.CartProduct
	context := helper.GetContextData(request)
	if context == nil {
		logrus.Error("AddProductToCart: Error in getting context data")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&cartItem)
	if err != nil {
		logrus.Error("AddProductToCart: Error in decoding json %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	//
	//productID, err := uuid.FromString(cartItem.ProductId)
	//if err != nil {
	//	logrus.Error(" AddProductToCart:Error in conversion: %v", err)
	//	writer.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	inventoryId, err := helper.GetInventoryId(cartItem.ProductId)
	if err != nil {
		logrus.Error(" AddProductToCart:Error in getting inventory Id: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	productQuantity, err := helper.GetCartItemInventoryQuantity(inventoryId)
	if err != nil {
		logrus.Error(" AddProductToCart:Error in getting product quantity: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if cartItem.Quantity > productQuantity {
		logrus.Error(" AddProductToCart: Please enter a valid quantity %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	//sessionId, err := uuid.FromString(context.SessionID)
	//if err != nil {
	//	logrus.Error(" AddProductToCart:Error in conversion: %v", err)
	//	writer.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//userId,err:=uuid.FromString(context.ID)
	//if err!=nil{
	//	logrus.Error(" AddProductToCart:Error in conversion: %v", err)
	//	writer.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	sessionID, err := helper.CreateCartSession(context.ID)
	sessionId, err := uuid.FromString(sessionID)
	if err != nil {
		logrus.Error(" AddProductToCart:Error in conversion: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = helper.AddToCart(sessionId, cartItem)
	if err != nil {
		logrus.Error(" AddProductToCart: Error in adding item to cart %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}

func AddUserAddress(writer http.ResponseWriter, request *http.Request) {
	var userAddress models.UserAddress
	err := json.NewDecoder(request.Body).Decode(&userAddress)
	if err != nil {
		logrus.Error("AddUserAddress: Error in decoding json %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	context := helper.GetContextData(request)
	if context == nil {
		logrus.Error("AddProductToCart: Error in getting context data")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := uuid.FromString(context.ID)
	if err != nil {
		logrus.Error("AddProductToCart: Error in conversion")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = helper.AddUserAddress(userId, userAddress)
	if err != nil {
		logrus.Error("AddProductToCart: Error in adding user address")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

}

func UpdateUserAddress(writer http.ResponseWriter, request *http.Request) {
	var userAddress models.UserAddress
	addressID := chi.URLParam(request, "addressID")
	err := json.NewDecoder(request.Body).Decode(&userAddress)
	if err != nil {
		logrus.Error("UpdateUserAddress: Error in decoding json %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	addressId, err := uuid.FromString(addressID)
	if err != nil {
		logrus.Error("Error in conversion: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = helper.UpdateAddress(addressId, userAddress)
	if err != nil {
		logrus.Error("UpdateUserAddress: Error in updating UserAddress %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

}

func DeleteUserAddress(writer http.ResponseWriter, request *http.Request) {
	addressID := chi.URLParam(request, "addressID")
	addressId, err := uuid.FromString(addressID)
	if err != nil {
		logrus.Error("Error in conversion: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = helper.DeleteAddress(addressId)
	if err != nil {
		logrus.Error("DeleteUserAddress: Error in deleting UserAddress %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

}

func AddPaymentDetails(writer http.ResponseWriter, request *http.Request) {
	var paymentData models.PaymentDetails
	err := json.NewDecoder(request.Body).Decode(&paymentData)
	if err != nil {
		logrus.Error("AddPaymentDetails: Error in decoding json %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := helper.GetContextData(request)
	if ctx == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := uuid.FromString(ctx.ID)
	if err != nil {
		logrus.Error("Error in conversion: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = helper.AddPaymentDetails(userId, paymentData)
	if err != nil {
		logrus.Error("AddPaymentDetails: Error in adding details %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}

func BuyProduct(writer http.ResponseWriter, request *http.Request) {
	productID := chi.URLParam(request, "productID")
	productId, err := uuid.FromString(productID)
	if err != nil {
		logrus.Error("BuyProduct:Error in conversion: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := helper.GetContextData(request)
	if ctx == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	//sessionId, err := uuid.FromString(ctx.SessionID)
	//if err != nil {
	//	logrus.Error("BuyProduct:Error in conversion: %v", err)
	//	return
	//}

	//userId, err := uuid.FromString(ctx.ID)
	//if err != nil {
	//	logrus.Error("BuyProduct:Error in conversion: %v", err)
	//	return
	//}

	CartItemQuantity, err := helper.GetProductQuantity(productId)
	if err != nil {
		logrus.Error("BuyProduct:Error in fetching cart quantity: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	inventoryId, err := helper.GetInventoryId(productId)
	if err != nil {
		logrus.Error("BuyProduct:Error in getting inventory Id: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	productQuantity, err := helper.GetCartItemInventoryQuantity(inventoryId)
	if err != nil {
		logrus.Error("BuyProduct:Error in getting product quantity: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionId, err := helper.GetCartSessionId(productId)
	if err != nil {
		logrus.Error("BuyProduct:Error in conversion: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	RemainingQuantity := productQuantity - CartItemQuantity
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err = helper.UpdateInventoryQuantity(inventoryId, RemainingQuantity)
		if err != nil {
			logrus.Error("BuyProduct:Error Updating product quantity: %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return err
		}

		err = helper.DeleteCart(sessionId, tx)
		if err != nil {
			logrus.Error("BuyProduct:Error in deleting cart : %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return err
		}

		return err

	})

	if txErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func GetOrderDetails(writer http.ResponseWriter, request *http.Request) {
	ctx := helper.GetContextData(request)
	if ctx == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := uuid.FromString(ctx.ID)
	if err != nil {
		logrus.Error("AddOrderDetails:Error in conversion: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	paymentID := chi.URLParam(request, "paymentID")
	paymentId, err := uuid.FromString(paymentID)
	if err != nil {
		logrus.Error("AddOrderDetails:Error in conversion: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	productID := chi.URLParam(request, "productID")
	productId, err := uuid.FromString(productID)
	if err != nil {
		logrus.Error("AddOrderDetails:Error in conversion: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	price, err := helper.GetProductPrice(productId)
	if err != nil {
		logrus.Error("AddOrderDetails:Error in getting product price: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	totalQuantity, err := helper.GetProductQuantity(productId)
	if err != nil {
		logrus.Error("AddOrderDetails:Error in getting product quantity: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	quantity := float64(totalQuantity)
	totalPrice := price * quantity
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err = helper.AddOrderDetails(userId, paymentId, totalPrice)
		if err != nil {
			logrus.Error("AddOrderDetails: Error in adding product details %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return err
		}
		orderDetails, err := helper.ShowOrderDetails(userId, tx)
		if err != nil {
			logrus.Error("AddOrderDetails: Error in showing product details %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return err
		}
		jsonData, jsonErr := json.Marshal(orderDetails)
		if jsonErr != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return err
		}

		_, err = writer.Write(jsonData)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return err

		}
		return err

	})
	if txErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func DeleteUserAccount(writer http.ResponseWriter, request *http.Request) {
	ctx := helper.GetContextData(request)
	if ctx == nil {
		logrus.Error("DeleteUserAccount:error in conversion")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := uuid.FromString(ctx.ID)
	if err != nil {
		logrus.Error("DeleteUserAccount:Error in conversion: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = helper.DeleteUserAccount(userId)
	if err != nil {
		logrus.Error("DeleteUserAccount:error in deleting user account")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

}
