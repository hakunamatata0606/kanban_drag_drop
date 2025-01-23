#include "datamodel_rep.h"

#include <QJsonDocument>
#include <QJsonParseError>
#include <QJsonArray>
#include <QJsonValue>
#include <QJsonObject>
#include <QDebug>
#include <QDebugStateSaver>

#include "datamodel.h"
#include "config.h"

DataModel_Rep::DataModel_Rep(DataModel *parent)
    : QObject{parent},
    m_state{State::Idle},
    m_wsk{QString(), QWebSocketProtocol::VersionLatest, this}
{
    QObject::connect(&m_wsk, &QWebSocket::connected, this, &DataModel_Rep::onConnected);
    QObject::connect(&m_wsk, &QWebSocket::disconnected, this, &DataModel_Rep::onDisconnected);
    QObject::connect(&m_wsk, &QWebSocket::errorOccurred, this, &DataModel_Rep::onErrorOccured);
    QObject::connect(&m_wsk, &QWebSocket::binaryMessageReceived, this, &DataModel_Rep::onBinaryMessageReceived);
    Config& config = Config::getInstance();
    m_wsk.open(QUrl{config.serverAddress()});
}

void DataModel_Rep::updateStatus(const QString &name, const QString &status)
{
    QByteArray msg;
    msg.push_back(static_cast<char>(Tag::UpdateTaskStatusMessage));
    QJsonObject obj;
    obj.insert("name", name);
    obj.insert("status", status);
    QJsonDocument doc(obj);
    msg.append(doc.toJson());

    m_wsk.sendBinaryMessage(msg);
}

void DataModel_Rep::onConnected()
{
    m_state = State::Connected;
    qDebug() << "DataModel_Rep::onConnected(): connected to server";
}

void DataModel_Rep::onDisconnected()
{
    m_state = State::Idle;
    qDebug() << "DataModel_Rep::onDisconnected(): disconnected from server";
}

void DataModel_Rep::onErrorOccured(QAbstractSocket::SocketError error)
{
    qDebug() << "DataModel_Rep::onErrorOccured(): got error: " << error;
}

void DataModel_Rep::onBinaryMessageReceived(const QByteArray &message)
{
    static const char* func = "DataModel_Rep::onBinaryMessageReceived(): ";
    QByteArray msg = message;
    if (!message.size()) {
        qDebug() << func << "receive empty message??";
        return;
    }
    uint8_t tag = static_cast<uint8_t>(*static_cast<const char*>(msg.data()));
    msg.removeFirst();
    switch (tag) {
    case Tag::ListStatusUpdate: {
        qDebug() << func << "receive list status update";
        handleListStatusUpdate(msg);
        break;
    }
    case Tag::ListTaskUpdate: {
        qDebug() << func << "receive list task update";
        handleTaskMapUpdate(msg);
        break;
    }
    default: {
        qDebug() << func << "receive unknown tag ??";
        break;
    }
    }
}

void DataModel_Rep::handleListStatusUpdate(const QByteArray &msg)
{
    static const char* func = "DataModel_Rep::handleListStatusUpdate(): ";
    QJsonParseError parseError;
    QJsonDocument jsonDoc = QJsonDocument::fromJson(msg, &parseError);
    if (parseError.error != QJsonParseError::NoError) {
        qDebug() << func << "failed to unmarsal list status";
        return;
    }
    if (!jsonDoc.isArray()) {
        qDebug() << func << "error is not an array";
        return;
    }
    listStatus_t listStatus;
    foreach (const QVariant& v, jsonDoc.array().toVariantList()) {
        listStatus.append(v.value<QString>());
    }
    emit listStatusUpdated(listStatus);
}

void DataModel_Rep::handleTaskMapUpdate(const QByteArray &msg)
{
    static const char* func = "DataModel_Rep::handleListTaskUpdate(): ";
    QJsonParseError parseError;
    QJsonDocument jsonDoc = QJsonDocument::fromJson(msg, &parseError);
    if (parseError.error != QJsonParseError::NoError) {
        qDebug() << func << "failed to unmarsal list task";
        return;
    }
    if (!jsonDoc.isArray()) {
        qDebug() << func << "error is not an array";
        return;
    }
    taskMap_t taskMap;
    foreach (const QJsonValue& v, jsonDoc.array()) {
        QJsonObject obj = v.toObject();
        task_t t;
        //expect to crash if there's error
        t.m_name = obj["name"].toString();
        t.m_title = obj["title"].toString();
        t.m_description = obj["description"].toString();
        t.m_status = obj["status"].toString();

        taskMap_t::iterator it = taskMap.find(t.m_status);
        if (it == taskMap.end()) {
            it = taskMap.insert(t.m_status, {});
        }
        it->push_back(std::move(t));
    }
    emit taskMapUpdated(taskMap);
}
