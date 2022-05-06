package main

func ExtractId(obj *JsonObject) string {
	return obj.GetString("_doc")
}

func ExtractTitle(obj *JsonObject) string {
	return obj.GetString("title")
}
