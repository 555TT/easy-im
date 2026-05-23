// easy-chat MongoDB 初始化脚本
// 适用服务：im-rpc、im-ws、task-mq
// 数据库名：easy-chat（连字符，区别于 MySQL 的 easy_chat）
// 认证：连接 URI 使用 root/easy-chat，默认 authSource=admin

const dbName = "easy-chat";
const target = db.getSiblingDB(dbName);

// ------------------------------------------------------------
// 集合（_id 全部使用默认 ObjectID；无 capped、无 validator、无 TTL）
// ------------------------------------------------------------

const collections = ["chat_log", "conversation", "conversations"];
const existing = new Set(target.getCollectionNames());

collections.forEach(function (name) {
    if (!existing.has(name)) {
        target.createCollection(name);
        print("created collection: " + name);
    } else {
        print("collection exists: " + name);
    }
});

// ------------------------------------------------------------
// 索引
// chat_log:        按会话拉历史消息（ListBySendTime）
// conversation:    全部读写都按 conversationId 命中，业务上要求唯一
// conversations:   每个用户一条文档，userId 唯一
// ------------------------------------------------------------

target.chat_log.createIndex(
    {conversationId: 1, sendTime: -1},
    {name: "idx_chatlog_conv_time"}
);

target.conversation.createIndex(
    {conversationId: 1},
    {name: "uk_conversation_conversationId", unique: true}
);

target.conversations.createIndex(
    {userId: 1},
    {name: "uk_conversations_userId", unique: true}
);

print("indexes ready on db " + dbName);


//docker exec -i mongo mongo \
//     --username root --password easy-chat --authenticationDatabase admin \
//     < go-im/deploy/mongo/init.js
