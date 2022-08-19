const axios = require("axios");
const getIPAddress = require("./getIpAddress");
const ipAddress = getIPAddress();
const register = async (address, port, name, id, tags) => {
    // 这里地址根据你部署的consul的来
    const url = "http://localhost:8500/v1/agent/service/register";
    const data = {
        Address: address,
        Port: port,
        Name: name,
        ID: id,
        Tags: tags,
        Check: {
            // 刚刚启动的express服务，但是这个地方地址不能写localhost
            HTTP: `http://${ipAddress}:3000/health`,
            Timeout: "5s",
            Interval: "5s",
            DeregisterCriticalServiceAfter: "10s",
        },
    };
    // 注册服务
    await axios.put(url, data);
};

module.exports = register;
