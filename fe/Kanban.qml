import QtQuick
import QtQuick.Controls
import fe

Item {
    id: root
    property DataModel dataModel

    Rectangle {
        id: rect
        anchors.fill: parent
        color: "#99ccff"
        GridView {
            id: grid
            anchors.fill: parent
            cellHeight: parent.height
            cellWidth: rect.width / elems.numOfStatus

            model: elems.listStatus
            delegate: Rectangle {
                id: drect
                width: grid.cellWidth
                height: grid.cellHeight
                border.color : "black"
                color: rect.color
                property int indexRect: index
                Column {
                    id: col
                    anchors.fill: parent
                    spacing: 10
                    Text {
                        id: title
                        text: status
                        anchors.horizontalCenter: parent.horizontalCenter
                        font.underline: true
                    }
                    Rectangle {
                        id: tasks
                        width: col.width
                        height: col.height - col.spacing - title.height

                        border.color: "black"
                        color: "green"

                        ListView {
                            id: listView
                            anchors.fill: parent
                            spacing: 10
                            model: elems.listTasks[drect.indexRect]
                            delegate: Card {
                                id: card
                                width: tasks.width
                                height: 100
                                name: _name
                                title: _title
                                status: _status
                                color: "orange"
                                onCardStatusChanged: function (n, s, change) {
                                    elems.changeStatus(n, s, change);
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    Item {
        id: elems
        property ListModel listStatus: ListModel{}
        property list<ListModel> listTasks: [];
        property int numOfStatus: 1
        Connections {
            target: root.dataModel
            function onListStatusChanged() {
                var listStatus = root.dataModel.listStatus;
                elems.listStatus.clear();
                elems.numOfStatus = listStatus.length;
                listStatus.forEach((elem) => elems.listStatus.append({status: elem}));
            }

            function onTaskMapChanged() {
                elems.listTasks = [];
                for (var i = 0; i < elems.listStatus.rowCount(); i++) {
                    var status = elems.listStatus.get(i).status;
                    elems.listTasks.push(mod.createObject());
                    var tasks = root.dataModel.getListTaskWithStatus(status);
                    tasks.forEach((e) => elems.listTasks[i].append({_name: e.name, _title: e.title, _status: e.status}));
                }
                console.log(elems.listTasks[4].rowCount());
            }
        }

        function changeStatus(name, status, change) {
            var nextStatus = getNextStatus(status, change);
            if (nextStatus === status) {
                return;
            }
            root.dataModel.updateStatus(name, nextStatus);
        }

        function getNextStatus(status, change) {
            for (var i = 0; i < elems.listStatus.rowCount(); i++) {
                if (elems.listStatus.get(i).status === status) {
                    var nextIdx = i + change;
                    if (nextIdx < 0 || nextIdx >= elems.listStatus.rowCount()) {
                        return status;
                    }
                    return elems.listStatus.get(nextIdx).status
                }
            }
            console.assert(0);
        }
    }

    Component {
        id: mod
        ListModel {}
    }
}
