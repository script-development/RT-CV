class CV {
  static example() {
    return CV(
      referenceNumber: "d1b77757-144a-4e67-871d-3a3c0f910c9b",
      createdAt: DateTime.now(),
      lastChanged: DateTime.now(),
      personalDetails: PersonalDetails(
        firstName: "Meneer",
        surName: "banaanmans",
        email: "some-email@example.com",
        phoneNumber: "06 11111111",
        dob: DateTime.now(),
        streetName: "Straatnaam",
        houseNumber: "15",
        houseNumberSuffix: "a",
        zip: "1234AB",
        city: "Groningen",
        country: "Nederland",
      ),
      workExperiences: [
        WorkExperience(
          profession: "Gangster",
          employer: "Gangsters B.V.",
          startDate: DateTime.utc(2000, 1, 1),
          endDate: DateTime.utc(2010, 1, 1),
          description: "Doing that gangster shit you know",
        ),
        WorkExperience(
          profession: "Bel meneer",
          employer: "De bel mensen",
          startDate: DateTime.utc(1995, 2, 1),
          endDate: DateTime.utc(1998, 5, 1),
          description:
              "Je weet wel die hinderlijke mensen die je opbellen met onzin waar je helemaal geen belang bij hebt",
        ),
      ],
      educations: [
        Education(
          "HBO",
          institute: "Hogeschool",
          startDate: DateTime.utc(2000, 1, 1),
          endDate: DateTime.utc(2010, 1, 1),
        ),
        Education(
          "MBO",
          institute: "Some school",
          startDate: DateTime.utc(1995, 2, 1),
          endDate: DateTime.utc(1998, 5, 1),
        ),
      ],
      courses: [
        Course(
          "beeing smart",
          institute: "Hogeschool",
          startDate: DateTime.utc(2010, 1, 1),
          endDate: DateTime.utc(2011, 1, 1),
        ),
      ],
      driversLicenses: ["B", "C", "D"],
      languages: [
        Language(
          "Nederlands",
          LanguageLevel.good,
          LanguageLevel.excellent,
        ),
        Language(
          "Duits",
          LanguageLevel.reasonable,
          LanguageLevel.good,
        ),
      ],
    );
  }

  CV({
    required this.referenceNumber,
    this.createdAt,
    this.lastChanged,
    this.educations,
    this.courses,
    this.workExperiences,
    this.preferredJobs,
    this.languages,
    required this.personalDetails,
    this.driversLicenses,
  });

  final String referenceNumber;
  final DateTime? createdAt;
  final DateTime? lastChanged;
  final List<Education>? educations;
  final List<Course>? courses;
  final List<WorkExperience>? workExperiences;
  final List<String>? preferredJobs;
  final List<Language>? languages;
  final PersonalDetails personalDetails;
  final List<String>? driversLicenses;
}

class Education {
  Education(
    this.name, {
    this.description,
    this.institute,
    this.isCompleted,
    this.hasDiploma,
    this.startDate,
    this.endDate,
  });

  final String name;
  final String? description;
  final String? institute;
  final bool? isCompleted;
  final bool? hasDiploma;
  final DateTime? startDate;
  final DateTime? endDate;
}

class Course {
  Course(
    this.name, {
    this.institute,
    this.startDate,
    this.endDate,
    this.isCompleted,
    this.description,
  });

  final String name;
  final String? institute;
  final DateTime? startDate;
  final DateTime? endDate;
  final bool? isCompleted;
  final String? description;
}

class WorkExperience {
  WorkExperience({
    this.profession,
    this.description,
    this.startDate,
    this.endDate,
    this.stillEmployed,
    this.employer,
    this.weeklyHoursWorked,
  });

  final String? profession;
  final String? description;
  final DateTime? startDate;
  final DateTime? endDate;
  final bool? stillEmployed;
  final String? employer;
  final int? weeklyHoursWorked;
}

enum LanguageLevel {
  unknown,
  reasonable,
  good,
  excellent,
}

Map<LanguageLevel, String> humanLanguageLevel = {
  LanguageLevel.unknown: "Onbekend",
  LanguageLevel.reasonable: "Redelijk",
  LanguageLevel.good: "Goed",
  LanguageLevel.excellent: "Uitstekend",
};

Map<LanguageLevel, int> languageLevelToNr = {
  LanguageLevel.unknown: 0,
  LanguageLevel.reasonable: 1,
  LanguageLevel.good: 2,
  LanguageLevel.excellent: 3,
};

Map<int, LanguageLevel> nrToLanguageLevel = {
  0: LanguageLevel.unknown,
  1: LanguageLevel.reasonable,
  2: LanguageLevel.good,
  3: LanguageLevel.excellent,
};

const int maxLanguageLevelNr = 3;

class Language {
  Language(this.name, this.levelSpoken, this.levelWritten);

  final String name;
  final LanguageLevel levelSpoken;
  final LanguageLevel levelWritten;
}

class Competence {
  Competence(this.name, {this.description});

  String name;
  String? description;
}

class Interest {
  Interest(this.name, {this.description});

  String name;
  String? description;
}

class PersonalDetails {
  PersonalDetails({
    this.initials,
    this.firstName,
    this.surNamePrefix,
    this.surName,
    this.dob,
    this.gender,
    this.streetName,
    this.houseNumber,
    this.houseNumberSuffix,
    this.zip,
    this.city,
    this.country,
    this.phoneNumber,
    this.email,
  });

  final String? initials;
  final String? firstName;
  final String? surNamePrefix;
  final String? surName;
  final DateTime? dob;
  final String? gender;
  final String? streetName;
  final String? houseNumber;
  final String? houseNumberSuffix;
  final String? zip;
  final String? city;
  final String? country;
  final String? phoneNumber;
  final String? email;

  bool get hasAddress =>
      streetName != null && houseNumber != null && city != null;

  String get fullName {
    String resp = this.firstName ?? '';

    if (this.surName != null) {
      if (resp.isNotEmpty) resp += " ";

      resp += (this.surNamePrefix != null ? this.surNamePrefix! + " " : '') +
          this.surName!;
    }

    return resp;
  }
}
