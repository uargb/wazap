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
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (Object.hasOwnProperty.call(mod, k)) result[k] = mod[k];
    result["default"] = mod;
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
var sulla_hotfix_1 = require("sulla-hotfix");
var axios_1 = __importDefault(require("axios"));
var fs = __importStar(require("fs"));
var mime = __importStar(require("mime-types"));
var qs = require('qs');
var apiBase = process.argv[0];
function base64_encode(file) {
    var data = fs.readFileSync(file);
    return Buffer.from(data).toString('base64');
}
function buildUrl(phone, url) {
    return apiBase + '/bot/' + phone + '/' + url;
}
function sendCard(card, phone, client) {
    return __awaiter(this, void 0, void 0, function () {
        var mediaData, mediaName;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    mediaData = 'data:{mime};base64,{base64}';
                    mediaName = '';
                    if (card.Image) {
                        mediaName = card.Image;
                        mediaData = mediaData.replace('{mime}', mime.lookup(card.Image));
                        mediaData = mediaData.replace('{base64}', base64_encode('./public/' + card.ManagerID + '-' + card.Image));
                    }
                    else if (card.Video) {
                        mediaName = card.Video;
                        mediaData = mediaData.replace('{mime}', mime.lookup(card.Video));
                        mediaData = mediaData.replace('{base64}', base64_encode('./public/' + card.ManagerID + '-' + card.Video));
                    }
                    else if (card.Attachment) {
                        mediaName = card.Attachment;
                        mediaData = mediaData.replace('{mime}', mime.lookup(card.Attachment));
                        mediaData = mediaData.replace('{base64}', base64_encode('./public/' + card.ManagerID + '-' + card.Attachment));
                    }
                    if (!(mediaName.length > 0)) return [3 /*break*/, 2];
                    return [4 /*yield*/, client.sendFile(phone + '@c.us', mediaData, mediaName, card.Text)];
                case 1:
                    _a.sent();
                    return [3 /*break*/, 4];
                case 2: return [4 /*yield*/, client.sendText(phone + '@c.us', card.Text)];
                case 3:
                    _a.sent();
                    _a.label = 4;
                case 4: return [2 /*return*/];
            }
        });
    });
}
function mailingUpdates(client) {
    return __awaiter(this, void 0, void 0, function () {
        var response, card_1, error_1;
        var _this = this;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    _a.trys.push([0, 2, 3, 4]);
                    return [4 /*yield*/, axios_1.default.get(buildUrl("0", 'mailing'))];
                case 1:
                    response = _a.sent();
                    if (response.data.ok) {
                        card_1 = response.data.data.card;
                        response.data.data.phones.forEach(function (phone) { return __awaiter(_this, void 0, void 0, function () {
                            return __generator(this, function (_a) {
                                switch (_a.label) {
                                    case 0: return [4 /*yield*/, sendCard(card_1, phone + '@c.us', client)];
                                    case 1:
                                        _a.sent();
                                        return [2 /*return*/];
                                }
                            });
                        }); });
                    }
                    return [3 /*break*/, 4];
                case 2:
                    error_1 = _a.sent();
                    console.error(error_1);
                    return [3 /*break*/, 4];
                case 3:
                    setTimeout(mailingUpdates, 10000, client);
                    return [7 /*endfinally*/];
                case 4: return [2 /*return*/];
            }
        });
    });
}
function start(client) {
    return __awaiter(this, void 0, void 0, function () {
        var _this = this;
        return __generator(this, function (_a) {
            client.onMessage(function (message) { return __awaiter(_this, void 0, void 0, function () {
                var phone, response_1, error_2, error_3;
                var _this = this;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            phone = message.from.split('@')[0];
                            _a.label = 1;
                        case 1:
                            _a.trys.push([1, 8, , 9]);
                            if (message.body === undefined) {
                                console.log(message);
                            }
                            return [4 /*yield*/, axios_1.default.get(buildUrl(phone, 'answer?message=' + encodeURIComponent(message.body)))];
                        case 2:
                            response_1 = _a.sent();
                            if (!response_1.data.ok) return [3 /*break*/, 7];
                            if (!(response_1.data.did === 'registered')) return [3 /*break*/, 6];
                            _a.label = 3;
                        case 3:
                            _a.trys.push([3, 5, , 6]);
                            return [4 /*yield*/, axios_1.default.get(buildUrl(phone, 'rename?name=' + encodeURIComponent(message.sender.pushname)))];
                        case 4:
                            _a.sent();
                            return [3 /*break*/, 6];
                        case 5:
                            error_2 = _a.sent();
                            console.error(error_2);
                            return [3 /*break*/, 6];
                        case 6:
                            response_1.data.data.forEach(function (card) { return __awaiter(_this, void 0, void 0, function () {
                                var fields, data, field;
                                return __generator(this, function (_a) {
                                    switch (_a.label) {
                                        case 0: return [4 /*yield*/, sendCard(card, phone, client)];
                                        case 1:
                                            _a.sent();
                                            if (!card.NotifyManager) return [3 /*break*/, 3];
                                            fields = qs.parse(response_1.data.costumer.Fields);
                                            data = 'Имя: ' + response_1.data.costumer.Name + '\n';
                                            data += 'Телефон: ' + response_1.data.costumer.Phone + '\n';
                                            for (field in fields) {
                                                if (field.includes('old')) {
                                                    continue;
                                                }
                                                data += field + ': ' + fields[field] + '\n';
                                            }
                                            return [4 /*yield*/, sendCard(card, response_1.data.manager.Phone + '@c.us', client)];
                                        case 2:
                                            _a.sent();
                                            _a.label = 3;
                                        case 3: return [2 /*return*/];
                                    }
                                });
                            }); });
                            _a.label = 7;
                        case 7: return [3 /*break*/, 9];
                        case 8:
                            error_3 = _a.sent();
                            console.error(error_3);
                            return [3 /*break*/, 9];
                        case 9: return [2 /*return*/];
                    }
                });
            }); });
            return [2 /*return*/];
        });
    });
}
console.log('starting client with API on ' + apiBase);
sulla_hotfix_1.create().then(function (client) {
    mailingUpdates(client);
    start(client);
}).catch(function (error) { return console.error(error); });
//# sourceMappingURL=index.js.map