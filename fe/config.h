#ifndef CONFIG_H
#define CONFIG_H

#include <QString>

class Config
{
    Config();
public:
    static Config& getInstance();

    const QString& serverAddress() const;

private:
    void init();
private:
    QString m_serverAddress;
};

#endif // CONFIG_H
