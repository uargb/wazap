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
Object.defineProperty(exports, "__esModule", { value: true });
var sulla_hotfix_1 = require("sulla-hotfix");
var mysql = require("promise-mysql");
var fs = require("fs");
var mime = require("mime-types");
var dbConfig = {
    host: 'wazap.cvtuyrclurh0.ap-south-1.rds.amazonaws.com',
    post: 3306,
    user: 'admin',
    password: '952368741',
    database: 'gym'
};
function translateQuery(query) {
    return query;
}
function base64_encode(file) {
    var data = fs.readFileSync(file);
    return Buffer.from(data).toString('base64');
}
function start(db, client) {
    return __awaiter(this, void 0, void 0, function () {
        var _this = this;
        return __generator(this, function (_a) {
            client.onMessage(function (message) { return __awaiter(_this, void 0, void 0, function () {
                var phone, result, managers, manager_1, response_1, qa_1, user, manager, qa, selected, unknown;
                var _this = this;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            phone = message.chatId.split('@')[0];
                            return [4 /*yield*/, db.query('select * from costumers where phone = ?', [phone])];
                        case 1:
                            result = _a.sent();
                            if (!(result.length == 0)) return [3 /*break*/, 6];
                            return [4 /*yield*/, db.query('select id, name, greeting from managers')
                                // по умолчанию скидываем админу
                            ];
                        case 2:
                            managers = _a.sent();
                            return [4 /*yield*/, db.query('select id, greeting, qa from managers where id = 1')];
                        case 3:
                            manager_1 = (_a.sent())[0];
                            managers.forEach(function (m) { return __awaiter(_this, void 0, void 0, function () {
                                return __generator(this, function (_a) {
                                    if (message.body.includes(manager_1.name)) {
                                        manager_1 = m; // если в сообщении есть имя менеджера - кинем ему
                                    }
                                    return [2 /*return*/];
                                });
                            }); });
                            return [4 /*yield*/, db.query('insert into costumers (phone, name, managerId) values(?, ?, ?)', [phone, message.sender.pushname, manager_1.id])];
                        case 4:
                            _a.sent();
                            response_1 = manager_1.greeting + "\n\n";
                            qa_1 = JSON.parse(manager_1.qa);
                            qa_1.forEach(function (item) {
                                if (item.show > 0) {
                                    response_1 += translateQuery(item.query) + ' - ' + item.description + "\n";
                                }
                            });
                            return [4 /*yield*/, client.sendText(message.from, response_1)];
                        case 5:
                            _a.sent();
                            return [2 /*return*/];
                        case 6:
                            user = result[0];
                            return [4 /*yield*/, db.query('select * from managers where id = ?', user.managerId)];
                        case 7:
                            manager = (_a.sent())[0];
                            qa = JSON.parse(manager.qa);
                            selected = [] // все карточки по выбранному запросу
                            ;
                            unknown = undefined // карточка "неопознанного" запроса
                            ;
                            qa.forEach(function (item) {
                                if (item.query === '<неизвестный>')
                                    unknown = item;
                                if (item.query.includes(message.body.trim()))
                                    selected.push(item);
                            });
                            if (selected.length == 0 && unknown)
                                selected.push(unknown);
                            selected.forEach(function (item) { return __awaiter(_this, void 0, void 0, function () {
                                var mediaData, mediaName;
                                return __generator(this, function (_a) {
                                    mediaData = 'data:{mime};base64,{base64}';
                                    mediaName = '';
                                    if (item.image) {
                                        mediaName = item.image;
                                        mediaData = mediaData.replace('{mime}', mime.lookup(item.image));
                                        mediaData = mediaData.replace('{base64}', base64_encode('./public/' + item.image));
                                    }
                                    else if (item.video) {
                                        mediaName = item.video;
                                        mediaData = mediaData.replace('{mime}', 'application/octet-stream');
                                        mediaData = mediaData.replace('{base64}', base64_encode('./public/' + item.video));
                                    }
                                    else if (item.attachment) {
                                        mediaName = item.attachment;
                                        mediaData = mediaData.replace('{mime}', mime.lookup(item.attachment));
                                        mediaData = mediaData.replace('{base64}', base64_encode('./public/' + item.attachment));
                                    }
                                    if (mediaName.length > 0)
                                        client.sendImage(message.from, mediaData, mediaName, item.text);
                                    else
                                        client.sendText(message.from, item.text);
                                    return [2 /*return*/];
                                });
                            }); });
                            return [2 /*return*/];
                    }
                });
            }); });
            return [2 /*return*/];
        });
    });
}
mysql.createConnection(dbConfig).then(function (db) {
    sulla_hotfix_1.create().then(function (client) {
        start(db, client);
    }).catch(function (error) { throw error; });
}).catch(function (error) { return console.error(error); });
//# sourceMappingURL=index.js.map