package Tables

import "Server/DataBase"

var AppData = DataBase.NewJsonTable[DataBase.AppData]("AppData")
