import QtQuick
import QtQuick.Controls

Item {
    id: root
    property string color
    property string name
    property string title
    property string status

    signal cardStatusChanged(string name, string status, int change);
    Rectangle {
        id: rect
        color: root.color
        anchors.fill: parent

        Column {
            id: col
            spacing: 10
            anchors.fill: parent
            Text {
                text: root.name
            }
            Text {
                id: tt
                text: root.title
            }
            Row {
                id: row
                Button {
                    text: "prev"
                    onClicked: {
                        root.cardStatusChanged(root.name, root.status, -1);
                    }
                }
                Button {
                    text: "next"
                    onClicked: {
                        root.cardStatusChanged(root.name, root.status, 1);
                    }
                }
            }
        }
    }
}
