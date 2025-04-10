import axios from 'axios';

const BASE_API_URL = process.env.REACT_APP_API_URL || 'http://localhost:5000';

interface Course {
  id?: number;
  name: string;
  description: string;
}

class DataCache<T> {
  private cacheStorage: Map<string, T> = new Map();

  retrieve(key: string): T | undefined {
    return this.cacheStorage.get(key);
  }

  store(key: string, value: T): void {
    this.cacheStorage.set(key, value);
  }
}

const courseDataCache = new DataCache<Course[]>();

async function getCoursesFromAPI(): Promise<Course[]> {
  const cacheKey = 'courses';
  const cachedCourses = courseDataCache.retrieve(cacheKey);
  if (cachedCourses) {
    return cachedCourses;
  }

  try {
    const response = await axios.get(`${BASE_API_URL}/courses`);
    courseDataCache.store(cacheKey, response.data);
    return response.data;
  } catch (error) {
    console.error('Error fetching courses:', error);
    throw error;
  }
}

async function addNewCourse(courseInfo: Course): Promise<Course> {
  try {
    const response = await axios.post(`${BASE_API_URL}/courses`, courseInfo);
    courseDataCache.store('courses', []);
    return response.data;
  } catch (error) {
    console.error('Error creating course:', error);
    throw error;
  }
}

async function modifyCourse(courseId: number, courseInfo: Course): Promise<Course> {
  try {
    const response = await axios.put(`${BASE_API_URL}/courses/${courseId}`, courseInfo);
    courseDataCache.store('courses', []);
    return response.data;
  } catch (error) {
    console.error('Error updating course:', error);
    throw error;
  }
}

async function removeCourse(courseId: number): Promise<void> {
  try {
    await axios.delete(`${BASE_API_URL}/courses/${courseId}`);
    courseDataCache.store('courses', []);
  } catch (error) {
    console.error('Error deleting course:', error);
    throw error;
  }
}

async function renderCourses() {
  const courses = await getCoursesFromAPI();
  const courseListElement = document.getElementById('courses-container');
  if (courseListElement) {
    courseListElement.innerHTML = '';
    courses.forEach(course => {
      const courseDetailElement = document.createElement('div');
      courseDetailElement.innerText = `Name: ${course.name}, Description: ${course.description}`;
      courseListElement.appendChild(courseDetailElement);
    });
  }
}

function initializeEventListeners() {
  const creationForm = document.getElementById('create-course-form');
  if (creationForm) {
    creationForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const nameField = (document.getElementById('course-name') as HTMLInputElement).value;
      const descriptionField = (document.getElementById('course-description') as HTMLInputElement).value;
      await addNewCourse({ name: nameField, description: descriptionField });
      await renderCourses();
    });
  }
}

async function initializeApplication() {
  await renderCourses();
  initializeEventListeners();
}

initializeApplication().catch(console.error);