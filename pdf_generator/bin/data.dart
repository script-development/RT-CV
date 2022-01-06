class CV {
  static example() {
    return CV(
      name: "Meneer banaanmans",
      reference: "d1b77757-144a-4e67-871d-3a3c0f910c9b",
      email: "some-email@example.com",
      phoneNr: "06 11111111",
      workExpr: [
        WorkExp(
          "Gangster",
          "Gangsters B.V.",
          DateTime.utc(2000, 1, 1),
          DateTime.utc(2010, 1, 1),
          "Doing that gangster shit you know",
        ),
        WorkExp(
          "Bel meneer",
          "De bel mensen",
          DateTime.utc(1995, 2, 1),
          DateTime.utc(1998, 5, 1),
          "Je weet wel die hinderlijke mensen die je opbellen met onzin waar je helemaal geen belang bij hebt",
        ),
      ],
      education: [
        Education(
          "HBO",
          "Hogeschool",
          DateTime.utc(2000, 1, 1),
          DateTime.utc(2010, 1, 1),
        ),
        Education(
          "MBO",
          "Some school",
          DateTime.utc(1995, 2, 1),
          DateTime.utc(1998, 5, 1),
        ),
      ],
      courses: [
        Education(
          "beeing smart",
          "Hogeschool",
          DateTime.utc(2010, 1, 1),
          DateTime.utc(2011, 1, 1),
        ),
      ],
      driversLicenses: ["B", "C", "D"],
      languageSkills: [
        LanguageSkill(
          "Nederlands",
          LanguageSkillLevel.excellent,
          LanguageSkillLevel.good,
        ),
        LanguageSkill(
          "Duits",
          LanguageSkillLevel.good,
          LanguageSkillLevel.reasonable,
        ),
      ],
    );
  }

  CV({
    required this.name,
    required this.reference,
    required this.email,
    required this.phoneNr,
    required this.workExpr,
    required this.education,
    required this.courses,
    required this.languageSkills,
    required this.driversLicenses,
  });

  final String name;
  final String reference;
  final String email;
  final String phoneNr;
  final List<WorkExp> workExpr;
  final List<Education> education;
  final List<Education> courses;
  final List<LanguageSkill> languageSkills;
  final List<String> driversLicenses;
}

enum LanguageSkillLevel {
  unknown,
  reasonable,
  good,
  excellent,
}

Map<LanguageSkillLevel, String> humanLanguageSkillLevel = {
  LanguageSkillLevel.unknown: "Onbekend",
  LanguageSkillLevel.reasonable: "Redelijk",
  LanguageSkillLevel.good: "Goed",
  LanguageSkillLevel.excellent: "Uitstekend",
};

class LanguageSkill {
  const LanguageSkill(this.name, this.reading, this.writing);

  final String name;
  final LanguageSkillLevel reading;
  final LanguageSkillLevel writing;
}

class Education {
  const Education(this.name, this.org, this.from, this.to);

  final String name;
  final String org;
  final DateTime from;
  final DateTime to;
}

class WorkExp {
  const WorkExp(
    this.name,
    this.company,
    this.from,
    this.to,
    this.description,
  );

  final String name;
  final String company;
  final DateTime from;
  final DateTime to;
  final String description;
}
