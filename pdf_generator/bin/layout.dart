import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';

class ListWithHeader {
  ListWithHeader(this.title, this.widgets);

  final String title;
  final List<Widget> widgets;

  get length => widgets.length;
}

class WrapLayoutBlock extends StatelessWidget {
  WrapLayoutBlock(this.listWithHeader);

  final ListWithHeader listWithHeader;

  @override
  Widget build(Context context) {
    List<List<Widget>> columns = [];
    for (var i = 0; i < listWithHeader.widgets.length; i++) {
      int idxInRow = i % 3;
      if (idxInRow == 0) {
        columns.add([]);
      }
      columns[columns.length - 1].add(Expanded(
        child: Padding(
          padding: EdgeInsets.only(
            left: idxInRow == 1 ? 4 : 0,
            right: idxInRow == 1
                ? 4
                : 0, // Only add horizontal padding to the second column
            top: i > 2
                ? 4
                : 0, // Only add padding to top if we are on the second row or more
          ),
          child: listWithHeader.widgets[i],
        ),
      ));
    }

    return LayoutBlockBase(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          ListTitle(listWithHeader.title),
          Column(
            children: columns
                .map((row) => Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: row,
                    ))
                .toList(),
          ),
        ],
      ),
    );
  }
}

class ColumnsLayoutBlock extends StatelessWidget {
  ColumnsLayoutBlock(this.widgets);

  final List<ListWithHeader> widgets;

  @override
  Widget build(Context context) {
    List<Widget> left = [];
    List<Widget> right = [];

    for (int i = 0; i < widgets.length; i++) {
      ListWithHeader widget = widgets[i];
      ListTitle title = ListTitle(
        widget.title,
        margin: i > 1 ? EdgeInsets.only(top: 10) : null,
      );

      if (i % 2 == 0) {
        left.add(title);
        left.addAll(widget.widgets);
      } else {
        right.add(title);
        right.addAll(widget.widgets);
      }
    }

    return LayoutBlockBase(
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Padding(
              padding: EdgeInsets.only(right: PdfPageFormat.cm / 2),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: left,
              ),
            ),
          ),
          Expanded(
            child: Padding(
              padding: EdgeInsets.only(left: PdfPageFormat.cm / 2),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: right,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class LayoutBlockBase extends StatelessWidget {
  LayoutBlockBase({required this.child});

  final Widget child;

  @override
  Widget build(Context contex) {
    return Padding(
      padding: EdgeInsets.only(
        left: PdfPageFormat.cm,
        right: PdfPageFormat.cm,
        top: PdfPageFormat.cm / 2,
      ),
      child: child,
    );
  }
}

class ListTitle extends StatelessWidget {
  ListTitle(this.title, {this.margin});

  final String title;
  final EdgeInsets? margin;

  @override
  Widget build(Context context) {
    PdfColor themeColor = PdfColor.fromInt(0xffffe004);

    return Container(
      margin: margin,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            color: themeColor,
            child: Padding(
              padding: EdgeInsets.symmetric(horizontal: 10, vertical: 5),
              child: Text(
                title,
                overflow: TextOverflow.clip,
                style: TextStyle(
                  fontSize: 10,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
          ),
          Container(
            width: double.infinity,
            height: 2,
            color: themeColor,
          ),
        ],
      ),
    );
  }
}
