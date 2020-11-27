package db

var Database = Db().Database("clipboard")
var ClipCollection = Database.Collection("clips")
var UserCollection = Database.Collection("users")