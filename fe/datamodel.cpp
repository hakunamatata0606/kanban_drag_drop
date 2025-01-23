#include "datamodel.h"

#include "datamodel_rep.h"


DataModel::DataModel(QObject *parent)
    : QObject{parent}, m_rep{nullptr}
{
    m_rep = new DataModel_Rep(this);
    QObject::connect(m_rep, &DataModel_Rep::listStatusUpdated, this, &DataModel::setListStatus);
    QObject::connect(m_rep, &DataModel_Rep::taskMapUpdated, this, &DataModel::setTaskMap);
}

DataModel::~DataModel()
{
    delete m_rep;
}

listStatus_t DataModel::listStatus() const
{
    return m_listStatus;
}

QDebug operator<<(QDebug db, const task_t& task)
{
    QDebugStateSaver saver(db);
    db << "(name: " << task.m_name << ", title: " << task.m_title << ", description: " << task.m_description << ", status: " << task.m_status << ")";
    return db;
}

taskMap_t DataModel::taskMap() const
{
    return m_taskMap;
}

void DataModel::updateStatus(const QString &name, const QString &status)
{
    m_rep->updateStatus(name, status);
}

QList<task_t> DataModel::getListTaskWithStatus(const QString &status)
{
    taskMap_t::iterator it = m_taskMap.find(status);
    if (it == m_taskMap.end()) {
        return {};
    }
    return *(it);
}



void DataModel::setListStatus(const listStatus_t &listStatus)
{
    m_listStatus = listStatus;
    emit listStatusChanged();
}

void DataModel::setTaskMap(const taskMap_t &taskMap)
{
    m_taskMap = taskMap;
    emit taskMapChanged();
}


