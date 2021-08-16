package messageCreator

import (
	"digi-bot/model"
	"fmt"
)

func CreateNormalPriceChangeMsg(product model.ProductDto, newPrice int, oldPrice int) string {
	output := createHeader(product).
		AddNewLine().
		Append("🔹").Append("قیمت: ")

	if oldPrice == 0 {
		output = output.
			Append(fmt.Sprintf("ناموجود😱🔒 -> %s", Number(newPrice).AddComma()))
	} else if newPrice == 0 {
		output = output.
			Append(fmt.Sprintf("%s -> ناموجود😱🔒", Number(oldPrice).AddComma()))
	} else {
		output = output.
			Append(fmt.Sprintf("%s -> %s", Number(oldPrice).AddComma(), Number(newPrice).AddComma()))
	}

	output = output.Append(createProductDetailMsg(product)).AddNewLine()

	return output.ToString()
}

func CreatePreviewMsg(product model.ProductDto) string {
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

	output = output.Append(createProductDetailMsg(product))

	output = output.
		AddNewLine().
		AddNewLine().
		Append("✅ کالا با موفقیت ذخیره شد. برای اضافه کردن کالای جدید کافی است فقط آدرس آن را وارد کنید")

	return output.ToString()
}

func CreateNotAvailableMsg(product model.ProductDto) string {
	output := createHeader(product).
		AddNewLine().
		Append("🔹").
		Append(String("ناموجود😱🔒").
			Bold().
			ToString())

	output = output.Append(createProductDetailMsg(product)).AddNewLine()

	return output.ToString()
}

func CreateDeleteProductSuccessfulMsg(product model.ProductDto) string {
	output := createHeader(product).AddNewLine().Append("✅ با موفقیت از لیست پاک شد").AddNewLine()
	return output.ToString()
}

func createHeader(product model.ProductDto) String {
	output := String(`🟣`).Append(product.Name).
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

	deleteAll := String("/deleteall").
		Bold().
		AddNewLine().
		Append("این دستور همه محصولات شما را پاک میکند").
		AddNewLine().
		ToString()

	return add + _delete + deleteAll
}

func createProductDetailMsg(product model.ProductDto) string {
	output := String("\n")
	if product.Desc1 != "" {
		output = output.
			Append("🔹").
			Append(product.Desc1).
			AddNewLine()
	}
	if product.Desc2 != "" {
		output = output.
			Append("🔹").
			Append(product.Desc2).
			AddNewLine()
	}
	if product.Desc3 != "" {
		output = output.
			Append("🔹").
			Append(product.Desc3).
			AddNewLine()
	}

	return output.ToString()
}

func CreateProductListMsg(products []string) string {
	str := String("")
	for _, product := range products {
		str = str.Append("🔹").
			Append(product).
			AddNewLine().
			AddNewLine()
	}
	return str.AddNewLine().AddNewLine().ToString()
}
