var createError = require("http-errors");
var express = require("express");
var path = require("path");
var cookieParser = require("cookie-parser");
var logger = require("morgan");

var indexRouter = require("./routes/index");
var usersRouter = require("./routes/users");
var register = require("./consul/register");

var app = express();

// view engine setup
app.set("views", path.join(__dirname, "views"));
app.set("view engine", "jade");

app.use(logger("dev"));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, "public")));

app.use("/", indexRouter);
app.use("/users", usersRouter);

app.get("/health", (req, res) => {
    res.json({ message: "成功" });
});
const getIPAddress = require("./consul/getIpAddress");
const ipAddress = getIPAddress();
app.get("/register", (req, res) => {
    const port = 50000;
    // 这里的name尽量和id保存一致，也可以不一致
    register(ipAddress, port, "node-consul", "node-consul", [
        "test",
        "node",
        "axios",
    ]);
    res.json({ message: `注册成功：${ipAddress}:${port}` });
});

// catch 404 and forward to error handler
app.use(function (req, res, next) {
    next(createError(404));
});

// error handler
app.use(function (err, req, res, next) {
    // set locals, only providing error in development
    res.locals.message = err.message;
    res.locals.error = req.app.get("env") === "development" ? err : {};

    // render the error page
    res.status(err.status || 500);
    res.render("error");
});

module.exports = app;
