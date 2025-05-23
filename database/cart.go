package database

import (
	"Ecommerce/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"time"
)

var (
	ErrCantFindProduct        = errors.New("can't find product")
	ErrCantDecodeProducts     = errors.New("can't find product")
	ErrUserIdIsNotValid       = errors.New("this user is not valid")
	ErrCantUpdateUser         = errors.New("can't add this product to the cart")
	ErrCantRemoveItemFromCart = errors.New("can't remove this item from the cart")
	ErrCantGetItem            = errors.New("unable to get item from cart")
	ErrCantBuyCartItem        = errors.New("can't update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {

	dbSearch, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser
	err = dbSearch.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{key: "userCart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}

	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D(primitive.E{Key: "_id", Value: id})
	update := bson.M{"$pull": bson.M{"userCart": bson.M{"_id": productID}}}

	_, err = UpdateMany(ctx, filter, update)

	if err != nil {
		return ErrCantRemoveItemFromCart
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	// fetch the cart of the user
	// find the cart total
	// create an order with the items
	// added order to the user Collection
	// added items in the cart to order list
	// empty up the cart
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getCartItems models.User
	var orderCart models.Order

	orderCart.OrderId = primitive.NewObjectID()
	orderCart.OrderedAt = time.Now()
	orderCart.OrderCart = make([]models.ProductUser, 0)
	orderCart.PaymentMethod.COD = true

	unWind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$userCart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$userCart.price"}}}}}}
	currentResult, err := userCollection.Aggregate(ctx, mongo.Pipeline{unWind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getUserCart []bson.M
	if err = currentResult.Decode(&getUserCart); err != nil {
		panic(err)
	}

	var totalPrice int32
	for _, user_item := range getUserCart {
		price := user_item["total"]
		totalPrice = price.(int32)
	}
	orderCart.Price = int(totalPrice)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.UserCart}}}

	_, err = userCollection.UpdateMany(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	userCartEmpty := make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "userCart", Value: userCartEmpty}}}}
	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		return ErrCantBuyCartItem
	}

	return nil
}

func InstanceBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var productDetails models.ProductUser
	var orderDetails models.Order

	orderDetails.OrderId = primitive.NewObjectID()
	orderDetails.OrderedAt = time.Now()
	orderDetails.OrderCart = make([]models.ProductUser, 0)
	orderDetails.PaymentMethod.COD = true

	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&productDetails)
	if err != nil {
		log.Println(err)
	}

	orderDetails.Price = productDetails.Price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": productDetails}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	return nil
}
