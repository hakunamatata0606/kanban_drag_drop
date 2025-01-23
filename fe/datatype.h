#ifndef DATATYPE_H
#define DATATYPE_H

#include <QString>
#include <QList>
#include <QMap>
#include <QObject>

typedef struct task {
    Q_GADGET
public:
    QString m_name;
    QString m_title;
    QString m_description;
    QString m_status;
    Q_PROPERTY(QString name MEMBER m_name)
    Q_PROPERTY(QString title MEMBER m_title)
    Q_PROPERTY(QString description MEMBER m_description)
    Q_PROPERTY(QString status MEMBER m_status)
}task_t;

enum Tag {
    CreateTaskMessage = 0,
    DeleteTaskMessage,
    UpdateTaskStatusMessage,
    GetListTask,
    MessageAck,
    ListStatusUpdate,
    ListTaskUpdate
};

typedef QMap<QString, QList<task_t>> taskMap_t;
typedef QList<QString> listStatus_t;

QDebug operator<<(QDebug db, const task_t& task);

#endif // DATATYPE_H
