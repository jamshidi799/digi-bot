package messageCreator

import (
	"digi-bot/model"
	"fmt"
)

func CreateNormalPriceChangeMsg(object model.Object, newPrice int, oldPrice int) string {
	output := createHeader(object).
		Append(fmt.Sprintf("%s -> %s", Number(oldPrice).AddComma(), Number(newPrice).AddComma()))
	return output.ToString()
}

func CreatePreviewMsg(object model.Object) string {
	output := createHeader(object)
	if object.Price != object.OldPrice {
		output.
			Append("قیمت: ").
			Append(fmt.Sprintf("%s -> %s",
				Number(object.OldPrice).
					AddComma().
					Strike(),
				Number(object.Price).
					AddComma()))
	} else {
		output.
			Append("قیمت: ").
			Append(Number(object.Price).AddComma().ToString())
	}

	output.
		AddNewLine().
		AddNewLine().
		Append("کالا با موفقیت ذخیره شد. برای اضافه کردن کالای جدید کافی است فقط آدرس آن را وارد کنید")

	return output.ToString()
}

func CreateNotAvailableMsg(obj model.Object) string {
	output := createHeader(obj).
		Append("ناموجود!")

	return output.ToString()
}

func createHeader(obj model.Object) String {
	output := String(obj.Name).
		Bold().
		ToLink(obj.Url).
		AddNewLine()

	return output
}
