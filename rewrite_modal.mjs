import fs from 'fs';
const c = fs.readFileSync('frontend-vue/src/pages/admin/TeachersPage.vue', 'utf8');

let result = c.replace(
  "const createForm = ref({ name: '', nickname: '', subject: '', grade_team: '', phone: '', email: '', password: '' })\nconst createAssignments = ref<{ class_id: number | null; role: Role; subject: string }[]>([])\nconst createLoading = ref(false)\nconst createGradeFilter = ref('')\nconst filteredCreateClasses = computed(() => {\n  if (!createGradeFilter.value) return classes.value\n  return classes.value.filter(c => c.grade === createGradeFilter.value)\n})",
  "const createForm = ref({ name: '', nickname: '', subject: '', grade_team: '', phone: '', email: '', password: '' })\nconst createAssignments = ref<{ class_id: number; class_name: string; subject: string }[]>([])\nconst createLoading = ref(false)\nconst pendingGrade = ref('')\nconst pendingClassId = ref(null)\nconst pendingSubject = ref('')\nconst gradeClasses = computed(() => classes.value.filter(c => c.grade === pendingGrade.value))"
);

result = result.replace(
  "function openCreateModal() {\n  createForm.value = { name: '', nickname: '', subject: '', grade_team: '', phone: '', email: '', password: '' }\n  // 默认列出所有班级供选择\n  createGradeFilter.value = ''\n  showCreateModal.value = true\n}",
  "function openCreateModal() {\n  createForm.value = { name: '', nickname: '', subject: '', grade_team: '', phone: '', email: '', password: '' }\n  createAssignments.value = []; pendingGrade.value = ''; pendingClassId.value = null; pendingSubject.value = ''\n  showCreateModal.value = true\n}"
);

result = result.replace(
  "function addCreateAssignment() { if (!createGradeFilter.value) { toast.show('请先选择具体年级', 'info'); return } createAssignments.value.push({ class_id: null, role: 'subject_teacher', subject: '' }) }\nfunction removeCreateAssignment(idx: number) { createAssignments.value.splice(idx, 1) }",
  "function addClassAssignment() {\n  if (!pendingClassId.value || !pendingGrade.value) { toast.show('请选择年级和班级', 'info'); return }\n  const cls = classes.value.find(c => c.id === pendingClassId.value)\n  if (!cls) return\n  if (createAssignments.value.some(a => a.class_id === cls.id)) { toast.show('该班级已添加', 'info'); return }\n  createAssignments.value.push({ class_id: cls.id, class_name: cls.name, subject: pendingSubject.value || '默认科目' })\n  pendingClassId.value = null; pendingSubject.value = ''\n}\nfunction removeClassAssignment(idx) { createAssignments.value.splice(idx, 1) }"
);

result = result.replace(
  "payload.assignments = createAssignments.value.map(a => ({ class_id: a.class_id, role: a.role, subject: a.subject || undefined }))",
  "payload.assignments = createAssignments.value.map(a => ({ class_id: a.class_id, role: 'subject_teacher', subject: a.subject || undefined }))"
);

fs.writeFileSync('frontend-vue/src/pages/admin/TeachersPage.vue', result);
console.log('Script section updated');
