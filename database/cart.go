package database

import (
	"Ecommerce/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
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

func BuyItemFromCart() {}

func InstanceBuyer() {}
