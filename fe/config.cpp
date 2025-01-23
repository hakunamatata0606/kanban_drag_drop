#include "config.h"

#include <QSettings>

static const char* DEFAUL_SERVER_URL = "ws://127.0.0.1:8080/ws";

Config::Config()
{
    init();
}

Config &Config::getInstance()
{
    static Config instance;
    return instance;
}

const QString &Config::serverAddress() const
{
    return m_serverAddress;
}

void Config::init()
{
    QSettings settings("./settings.ini", QSettings::IniFormat);
    settings.beginGroup("config");

    QString serverAddress = settings.value("server_address").toString();
    if (serverAddress.isEmpty()) {
        qDebug() << "Failed to read server address from config";
        serverAddress = DEFAUL_SERVER_URL;
    }
    m_serverAddress = serverAddress;
}
