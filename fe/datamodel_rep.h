#ifndef DATAMODEL_REP_H
#define DATAMODEL_REP_H

#include <QObject>
#include <QWebSocket>

#include "datatype.h"

class DataModel;

class DataModel_Rep : public QObject
{
    Q_OBJECT

    enum class State {
        Idle,
        Connected
    };

public:
    explicit DataModel_Rep(DataModel* parent);

    void updateStatus(const QString& name, const QString& status);
public slots:
    void onConnected();
    void onDisconnected();
    void onErrorOccured(QAbstractSocket::SocketError error);

    void onBinaryMessageReceived(const QByteArray &message);
signals:
    void listStatusUpdated(const listStatus_t& listStatus);
    void taskMapUpdated(const taskMap_t& taskMap);
private:
    void handleListStatusUpdate(const QByteArray& msg);

    void handleTaskMapUpdate(const QByteArray& msg);
private:
    State m_state;
    QWebSocket m_wsk;
};

#endif // DATAMODEL_REP_H
