

function table_generating_script(MainTeacherSubjectObject) {
    html = "<table><thead><tr>";
    teachersubjects = Object.keys(MainTeacherSubjectObject);
    for (var i = 0; i < teachersubjects.length; i++) { html += "<td>" + teachersubjects[i] + "</td><td>subject listing</td>"; }
    html += "</tr></thead>";
    
    teachersubjects_html = [];
    max_length = -1;
    for (var i = 0; i < teachersubjects.length; i++) {
        subjects = MainTeacherSubjectObject[teachersubjects[i]];
        subjects_html = [];
        for (var j = 0; j < subjects.length; j++) {
            subject = subjects[j];
            subjects_html = subjects_html.concat(subtable_generating_script(subject));
        }
        teachersubjects_html.push(subjects_html);
        
        if (max_length < subjects_html.length) max_length = subjects_html.length;
    }
    
    console.log(teachersubjects_html);
    
    for (var i = 0; i < max_length; i++) {
        html += "<tr>";
        row = "";
        for (var j = 0; j < teachersubjects_html.length; j++) {
            row += (teachersubjects_html[j].length > i) ? teachersubjects_html[j][i] : "<td></td><td></td>";
        }
        console.log(row)
        html += row
        html += "</tr>";
    }
    return html;
}

function subtable_generating_script(list) {
    "ex: ['Modern Biology','Mr. Econome','Ms. Banfield','Ms. Prabhu','Ms.Hua']";
    len = list.length;
    num = len;
    elem0 = "<td>" + list[0] + "</td><td>" + list[1] + "</td>";
    newtabularlist = [elem0]; //Lists will be compiled in the end and delimited by </tr><tr>;
    for (var i = 2; i < len; i++) {
        newtabularlist.push("<td></td><td>" + list[i] + "</td>");
    }
    newtabularlist.push("<td></td><td></td>")
    return newtabularlist;
}

function subject_listing_script(MainTeacherSubjectObject) {
    
    keys = Object.keys(MainTeacherSubjectObject);
    subjecthtml = "";
    for (var i = 0; i < keys.length; i++) {
        subjecthtml += "<h1>" + keys[i] + "</h1>"
        subject_data = MainTeacherSubjectObject[keys[i]]
        for (var j = 0; j < subject_data.length; j++) {
            subjecthtml += "<h3>" + subject_data[j][0] + "</h3>";
            subjecthtml += "<div id='topic" + subject_data[j][0].replace(' ', '') + "'></div>";
        }
    }
    return subjecthtml;
    
}