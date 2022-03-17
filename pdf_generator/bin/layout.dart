import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';

import 'utils.dart';
import 'args.dart';

class ListWithHeader {
  ListWithHeader(this.icon, this.title, this.widgets);

  final IconData icon;
  final String title;
  final List<Widget> widgets;

  get length => widgets.length;
}

// class BigListLayoutBlock extends StatelessWidget {
//   BigListLayoutBlock(this.listWithHeader, this.style);

//   final ListWithHeader listWithHeader;
//   final Style style;

//   @override
//   Widget build(Context context) {
//     return Column(
//       children: [
//         ListTitle(
//           listWithHeader.icon,
//           listWithHeader.title,
//           style,
//         ),
//         ...listWithHeader.widgets.map((w) => ).toList(),
//         Padding(
//           padding: EdgeInsets.only(
//             left: style.layoutStyle == LayoutStyle.style_3 ? 10 : 0,
//           ),
//           child: Column(
//             crossAxisAlignment: CrossAxisAlignment.start,
//             children: columnWidgets,
//           ),
//         ),
//       ],
//     );
//   }
// }

class WrapLayoutBlock extends StatelessWidget {
  WrapLayoutBlock(this.listWithHeader, this.style);

  final ListWithHeader listWithHeader;
  final Style style;

  @override
  Widget build(Context context) {
    List<List<Widget>> columns = [
      [],
      [],
      [],
    ];

    for (var i = 0; i < listWithHeader.widgets.length; i++) {
      if (i > 2) {
        columns[i % 3].add(Padding(
          padding: const EdgeInsets.only(top: 4),
          child: listWithHeader.widgets[i],
        ));
      } else {
        columns[i % 3].add(listWithHeader.widgets[i]);
      }
    }

    List<Widget> columnWidgets = [];
    for (var i = 0; i < columns.length; i++) {
      columnWidgets.add(Expanded(
        child: Padding(
          padding: EdgeInsets.only(
            left: i == 0 ? 0 : 2,
            right: i == 2 ? 0 : 2,
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: columns[i],
          ),
        ),
      ));
    }

    return LayoutBlockBase(
      child: Container(
        decoration: BoxDecoration(
          border: style.layoutStyle == LayoutStyle.style_3
              ? Border(
                  left: BorderSide(
                    color: style.subHeaderBackgroundColor,
                    width: 2,
                  ),
                )
              : null,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            ListTitle(
              listWithHeader.icon,
              listWithHeader.title,
              style,
            ),
            Padding(
              padding: EdgeInsets.only(
                  left: style.layoutStyle == LayoutStyle.style_3 ? 10 : 0),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: columnWidgets,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class ColumnsLayoutBlock extends StatelessWidget {
  ColumnsLayoutBlock(this.widgets, this.style);

  final List<ListWithHeader> widgets;
  final Style style;

  @override
  Widget build(Context context) {
    List<Widget> left = [];
    List<Widget> right = [];

    for (int i = 0; i < widgets.length; i++) {
      ListWithHeader widget = widgets[i];

      Widget listWithHeaderToAdd = ColumnLayoutBlock(
        widget,
        style,
        EdgeInsets.only(top: i > 1 ? 10 : 0),
      );

      if (i % 2 == 0) {
        left.add(listWithHeaderToAdd);
      } else {
        right.add(listWithHeaderToAdd);
      }
    }

    return LayoutBlockBase(
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Padding(
              padding: const EdgeInsets.only(right: PdfPageFormat.cm / 2),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: left,
              ),
            ),
          ),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.only(left: PdfPageFormat.cm / 2),
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

class ColumnLayoutBlock extends StatelessWidget {
  ColumnLayoutBlock(this.widget, this.style, [this.padding = EdgeInsets.zero]);

  final ListWithHeader widget;
  final Style style;
  final EdgeInsets padding;

  MightApplyStyle3(Widget child) {
    if (style.layoutStyle != LayoutStyle.style_3) return child;

    return Container(
      decoration: BoxDecoration(
        border: style.layoutStyle == LayoutStyle.style_3
            ? Border(
                left: BorderSide(
                  color: style.subHeaderBackgroundColor,
                  width: 2,
                ),
              )
            : null,
      ),
      child: child,
    );
  }

  @override
  Widget build(Context context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: padding,
          child: MightApplyStyle3(ListTitle(
            widget.icon,
            widget.title,
            style,
          )),
        ),
        ...widget.widgets
            .map((w) => Padding(
                padding: padding.copyWith(top: 0, bottom: 0),
                child: MightApplyStyle3(Padding(
                    padding: EdgeInsets.only(
                      left: style.layoutStyle == LayoutStyle.style_3 ? 10 : 0,
                    ),
                    child: w))))
            .toList(),
      ],
    );
  }
}

const layoutBlockBasePadding = EdgeInsets.only(
  left: PdfPageFormat.cm,
  right: PdfPageFormat.cm,
  top: PdfPageFormat.cm / 2,
);

class LayoutBlockBase extends StatelessWidget {
  LayoutBlockBase({required this.child});

  final Widget child;

  @override
  Widget build(Context contex) {
    return Padding(
      padding: layoutBlockBasePadding,
      child: child,
    );
  }
}

class ListTitle extends StatelessWidget {
  ListTitle(this.icon, this.title, this.style);

  final IconData icon;
  final String title;
  final Style style;

  @override
  Widget build(Context context) {
    var children = [
      Container(
        decoration: BoxDecoration(
          borderRadius: style.layoutStyle == LayoutStyle.style_2
              ? BorderRadius.circular(5)
              : null,
          color: style.subHeaderBackgroundColor,
        ),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Padding(
                padding: const EdgeInsets.only(right: 5),
                child: Icon(
                  icon,
                  size: 10,
                  color: style.subHeaderTextColor,
                ),
              ),
              Text(
                title,
                overflow: TextOverflow.clip,
                style: TextStyle(
                  color: style.subHeaderTextColor,
                  fontSize: 10,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
        ),
      ),
    ];

    if (style.layoutStyle == LayoutStyle.style_1) {
      children.add(Container(
        width: double.infinity,
        height: 2,
        color: style.subHeaderBackgroundColor,
      ));
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: children,
    );
  }
}
