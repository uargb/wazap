"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
exports.__esModule = true;
var sulla_hotfix_1 = require("sulla-hotfix");
var mysql = require("promise-mysql");
var dbConfig = {
    host: 'wazap.cvtuyrclurh0.ap-south-1.rds.amazonaws.com',
    post: 3306,
    user: 'admin',
    password: '952368741',
    database: 'gym'
};
function start(db, client) {
    var _this = this;
    client.onMessage(function (message) { return __awaiter(_this, void 0, void 0, function () {
        var phone, result, managers, newManager_1, response, user;
        var _this = this;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    phone = message.chatId.split('@')[0];
                    return [4 /*yield*/, db.query('select * from costumers where phone = ?', [phone])];
                case 1:
                    result = _a.sent();
                    if (!(result.length == 0)) return [3 /*break*/, 5];
                    return [4 /*yield*/, db.query('select id, name, greeting from managers')
                        // по умолчанию скидываем админу
                    ];
                case 2:
                    managers = _a.sent();
                    return [4 /*yield*/, db.query('select id, greeting from managers where id = 1')[0]];
                case 3:
                    newManager_1 = _a.sent();
                    managers.forEach(function (manager) { return __awaiter(_this, void 0, void 0, function () {
                        return __generator(this, function (_a) {
                            if (message.body.includes(manager.name)) {
                                newManager_1 = manager; // если в сообщении есть имя менеджера - кинем ему
                            }
                            return [2 /*return*/];
                        });
                    }); });
                    return [4 /*yield*/, db.query('insert into costumers (phone, name, managerId) values(?, ?, ?)', [phone, message.sender.pushname, newManager_1.id])];
                case 4:
                    _a.sent();
                    response = newManager_1.greeting;
                    // todo: организовать меню
                    client.sendText(message.from, response);
                    _a.label = 5;
                case 5:
                    user = result[0];
                    console.log(result);
                    return [2 /*return*/];
            }
        });
    }); });
}
mysql.createConnection(dbConfig).then(function (db) {
    sulla_hotfix_1.create().then(function (client) {
        start(db, client);
    })["catch"](function (error) { throw error; });
})["catch"](function (error) { return console.error(error); });
