import QtQuick
import QtQuick.Controls
import fe

Window {
    width: 640
    height: 480
    visible: true
    title: qsTr("Hello World")

    Kanban{
        id: root
        anchors.fill: parent
        dataModel: dataModel
    }


    DataModel {
        id: dataModel

        onListStatusChanged : {
            // console.log(listStatus);
            // listStatus.forEach((elem) => grid.items.append({name: elem}));
            // listStatusChanged(listStatus);
        }
        onTaskMapChanged: {
            // console.log("aloha: " + taskMap);
            // console.log("aloha1: " + getListTaskWithStatus("idea")[0]);
        }
    }
}
