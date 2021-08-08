package messageCreator

import (
	"digi-bot/model"
	"fmt"
)

func CreateNormalPriceChangeMsg(product model.Product, newPrice int, oldPrice int) string {
	output := createHeader(product)
	if oldPrice == 0 {
		output = output.
			Append(fmt.Sprintf("ناموجود😱🔒 -> %s", Number(oldPrice).AddComma()))
	} else {
		output = output.
			Append(fmt.Sprintf("%s -> %s", Number(oldPrice).AddComma(), Number(newPrice).AddComma()))
	}

	return output.ToString()
}

func CreatePreviewMsg(product model.Product) string {
	output := String(`🟣`).
		Append(product.Name).
		Bold().
		AddNewLine().
		AddNewLine().
		Append("🔹")

	if product.Price != product.OldPrice {
		output = output.
			Append("قیمت: ").
			Append(fmt.Sprintf("%s -> %s",
				Number(product.OldPrice).
					AddComma().
					Strike(),
				Number(product.Price).
					AddComma()))
	} else {
		if product.Price != 0 {
			output = output.
				Append("قیمت: ").
				Append(Number(product.Price).AddComma().ToString())
		} else {
			output = output.
				Append(String("ناموجود😱🔒").Bold().ToString())
		}
	}

	output = output.
		AddNewLine().
		Append("🔹").
		Append(product.Desc1).
		AddNewLine().
		Append("🔹").
		Append(product.Desc2)

	output = output.
		AddNewLine().
		AddNewLine().
		Append("✅ کالا با موفقیت ذخیره شد. برای اضافه کردن کالای جدید کافی است فقط آدرس آن را وارد کنید")

	return output.ToString()
}

func CreateNotAvailableMsg(product model.Product) string {
	output := createHeader(product).
		Append(String("ناموجود😱🔒").Bold().ToString())

	return output.ToString()
}

func CreateDeleteProductSuccessfulMsg(product model.Product) string {
	output := createHeader(product).Append("✅ با موفقیت از لیست پاک شد")
	return output.ToString()
}

func createHeader(product model.Product) String {
	output := String(product.Name).
		Bold().
		ToLink(product.Url).
		AddNewLine()

	return output
}

func CreateHelpMsg() string {
	add := String("+").
		Bold().
		Append("برای اضافه کردن, آدرس(url) محصول را وارد کنید").
		AddNewLine().
		AddNewLine().
		ToString()

	_delete := String("-").
		Bold().
		Append("برای حذف کردن فقط به کالای موردنظر ریپلای بزنید").
		AddNewLine().
		AddNewLine().
		ToString()

	deleteAll := String("/deleteAll").
		Bold().
		AddNewLine().
		Append("این دستور همه محصولات شما را پاک میکند").
		AddNewLine().
		ToString()

	return add + _delete + deleteAll
}
