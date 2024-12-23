import axios from 'axios';
const BASE_API_URL = process.env.REACT_APP_API_URL || 'http://localhost:5000';
interface Course {
    id?: number;
    name: string;
    description: string;
}
async function fetchCourses(): Promise<Course[]> {
    try {
        const response = await axios.get(`${BASE_API_URL}/courses`);
        return response.data;
    } catch (error) {
        console.error('Error fetching courses:', error);
        throw error;
    }
}
async function createCourse(courseData: Course): Promise<Course> {
    try {
        const response = await axios.post(`${BASE_API_URL}/courses`, courseData);
        return response.data;
    } catch (error) {
        console.error('Error creating course:', error);
        throw error;
    }
}
async function updateCourse(courseId: number, courseData: Course): Promise<Course> {
    try {
        const response = await axios.put(`${BASE_API_URL}/courses/${courseId}`, courseData);
        return response.data;
    } catch (error) {
        console.error('Error updating course:', error);
        throw error;
    }
}
async function deleteCourse(courseId: number): Promise<void> {
    try {
        await axios.delete(`${BASE_API_URL}/courses/${courseId}`);
    } catch (error) {
        console.error('Error deleting course:', error);
        throw error;
    }
}
async function displayCourses() {
    const courses = await fetchCourses();
    const coursesContainer = document.getElementById('courses-container');
    if (coursesContainer) {
        coursesContainer.innerHTML = ''; 
        courses.forEach(course => {
            const courseElement = document.createElement('div');
            courseElement.innerText = `Name: ${course.name}, Description: ${course.description}`;
            coursesContainer.appendChild(courseElement);
        });
    }
}
function setupEventListeners() {
    const createForm = document.getElementById('create-course-form');
    if (createForm) {
        createForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const nameInput = (document.getElementById('course-name') as HTMLInputElement).value;
            const descriptionInput = (document.getElementById('course-description') as HTMLInputElement).value;
            await createCourse({name: nameInput, description: descriptionInput});
            await displayCourses(); 
        });
    }
}
async function init() {
    await displayCourses();
    setupEventListeners();
}
init().catch(console.error);