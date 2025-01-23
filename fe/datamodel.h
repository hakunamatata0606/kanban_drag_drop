#ifndef DATAMODEL_H
#define DATAMODEL_H

#include <QObject>
#include <QQmlEngine>

#include "datatype.h"

class DataModel_Rep;

class DataModel : public QObject
{
    Q_OBJECT
    QML_ELEMENT
    Q_PROPERTY(listStatus_t listStatus READ listStatus WRITE setListStatus NOTIFY listStatusChanged FINAL)
    Q_PROPERTY(taskMap_t taskMap READ taskMap WRITE setTaskMap NOTIFY taskMapChanged FINAL)
public:
    explicit DataModel(QObject *parent = nullptr);

    ~DataModel();

    listStatus_t listStatus() const;

    taskMap_t taskMap() const;

    Q_INVOKABLE void updateStatus(const QString& name, const QString& status);

    Q_INVOKABLE QList<task_t> getListTaskWithStatus(const QString& status);
public slots:
    void setListStatus(const listStatus_t& listStatus);

    void setTaskMap(const taskMap_t& taskMap);
signals:
    void listStatusChanged();

    void taskMapChanged();
private:
    listStatus_t m_listStatus;
    taskMap_t m_taskMap;

    DataModel_Rep* m_rep;
};

#endif // DATAMODEL_H
