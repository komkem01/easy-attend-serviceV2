# Frontend API Documentation
## Easy Attend Service V2 - หน้าบ้าน

**เวอร์ชัน:** 2.0  
**วันที่:** 23 ตุลาคม 2025  
**Base URL:** `http://localhost:8080`

---

## 📋 สารบัญ

1. [Authentication APIs](#authentication-apis)
2. [Dashboard APIs](#dashboard-apis)
3. [Classroom Management](#classroom-management)
4. [Student Management](#student-management)
5. [Attendance Management](#attendance-management)
6. [วิธีการใช้งาน](#วิธีการใช้งาน)
7. [ตัวอย่างโค้ด](#ตัวอย่างโค้ด)

---

## 🔐 Authentication APIs

### 1. เข้าสู่ระบบ (Login)

**Endpoint:** `POST /api/v1/login`

**Request Body:**
```json
{
  "email": "teacher@example.com",
  "password": "password123"
}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "เข้าสู่ระบบสำเร็จ",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires": "2025-10-24T20:00:00Z",
  "data": {
    "id": 1,
    "email": "teacher@example.com",
    "first_name": "สมชาย",
    "last_name": "ใจดี",
    "phone": "0812345678",
    "school": {
      "id": 1,
      "name": "โรงเรียนสมชาย"
    },
    "prefix": {
      "id": 1,
      "name": "นาย"
    },
    "gender": {
      "id": 1,
      "name": "ชาย"
    }
  }
}
```

**Response (Error):**
```json
{
  "success": false,
  "message": "อีเมลหรือรหัสผ่านไม่ถูกต้อง",
  "error": "invalid email or password"
}
```

### 2. สมัครสมาชิก (Register)

**Endpoint:** `POST /api/v1/register`

**Request Body:**
```json
{
  "email": "newteacher@example.com",
  "password": "password123",
  "first_name": "สมหญิง",
  "last_name": "ใจดี",
  "phone": "0823456789",
  "school_name": "โรงเรียนใหม่",
  "prefix_id": 2,
  "gender_id": 2
}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "สมัครสมาชิกสำเร็จ",
  "data": {
    "id": 2,
    "email": "newteacher@example.com",
    "first_name": "สมหญิง",
    "last_name": "ใจดี",
    "phone": "0823456789"
  }
}
```

---

## 📊 Dashboard APIs

### 3. ข้อมูลหน้าแรก (Dashboard)

**Endpoint:** `GET /api/v1/dashboard`  
**Authentication:** Required (Bearer Token)

**Response:**
```json
{
  "success": true,
  "message": "ดึงข้อมูลสำเร็จ",
  "data": {
    "summary": {
      "total_classrooms": 5,
      "total_students": 150,
      "total_sessions": 25,
      "present_today": 142,
      "absent_today": 8
    },
    "classrooms": [
      {
        "id": 1,
        "name": "ม.1/1",
        "grade": "ม.1",
        "student_count": 30,
        "students": [
          {
            "id": 1,
            "student_no": "STD001",
            "first_name": "สมศรี",
            "last_name": "ใจดี",
            "prefix": {
              "id": 2,
              "name": "นางสาว"
            },
            "gender": {
              "id": 2,
              "name": "หญิง"
            }
          }
        ]
      }
    ],
    "recent_activities": [
      {
        "id": 1,
        "action": "เช็คชื่อ",
        "description": "เช็คชื่อห้อง ม.1/1 วันที่ 23/10/2025",
        "created_at": "2025-10-23T13:30:00Z"
      }
    ]
  }
}
```

---

## 🏫 Classroom Management

### 4. สร้างห้องเรียน

**Endpoint:** `POST /api/v1/classrooms`  
**Authentication:** Required

**Request Body:**
```json
{
  "name": "ม.2/3",
  "grade": "ม.2"
}
```

**Response:**
```json
{
  "success": true,
  "message": "สร้างห้องเรียนสำเร็จ",
  "data": {
    "id": 6,
    "name": "ม.2/3",
    "grade": "ม.2",
    "student_count": 0
  }
}
```

### 5. รายละเอียดห้องเรียน

**Endpoint:** `GET /api/v1/classrooms/{id}`  
**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "message": "ดึงข้อมูลสำเร็จ",
  "data": {
    "classroom": {
      "id": 1,
      "name": "ม.1/1",
      "grade": "ม.1",
      "student_count": 30
    },
    "attendance_stats": {
      "total_sessions": 20,
      "present_rate": 85,
      "absent_rate": 10,
      "late_rate": 5
    },
    "recent_sessions": [
      {
        "session_date": "2025-10-23",
        "total_students": 30,
        "present": 28,
        "absent": 1,
        "late": 1,
        "attendances": [
          {
            "student": {
              "id": 1,
              "student_no": "STD001",
              "first_name": "สมศรี",
              "last_name": "ใจดี"
            },
            "status": "present",
            "remark": ""
          }
        ]
      }
    ]
  }
}
```

---

## 👥 Student Management

### 6. เพิ่มนักเรียน

**Endpoint:** `POST /api/v1/students`  
**Authentication:** Required

**Request Body:**
```json
{
  "classroom_id": 1,
  "first_name": "สมหมาย",
  "last_name": "รักเรียน",
  "prefix_id": 1,
  "gender_id": 1
}
```

**Response:**
```json
{
  "success": true,
  "message": "เพิ่มนักเรียนสำเร็จ",
  "data": {
    "id": 31,
    "student_no": "STD031",
    "first_name": "สมหมาย",
    "last_name": "รักเรียน"
  }
}
```

### 7. รายการนักเรียนในห้อง

**Endpoint:** `GET /api/v1/classrooms/{id}/students`  
**Authentication:** Required

**Query Parameters:**
- `page` (optional): หน้าที่ต้องการ (default: 1)
- `limit` (optional): จำนวนต่อหน้า (default: 50, max: 100)

**Response:**
```json
{
  "success": true,
  "message": "ดึงข้อมูลสำเร็จ",
  "data": [
    {
      "id": 1,
      "student_no": "STD001",
      "first_name": "สมศรี",
      "last_name": "ใจดี",
      "prefix": {
        "id": 2,
        "name": "นางสาว"
      },
      "gender": {
        "id": 2,
        "name": "หญิง"
      }
    }
  ],
  "meta": {
    "page": 1,
    "limit": 50,
    "total": 30,
    "total_pages": 1
  }
}
```

---

## ✅ Attendance Management

### 8. บันทึกการเช็คชื่อ

**Endpoint:** `POST /api/v1/attendance`  
**Authentication:** Required

**Request Body:**
```json
{
  "classroom_id": 1,
  "session_date": "2025-10-23",
  "students": [
    {
      "student_id": 1,
      "status": "present",
      "remark": ""
    },
    {
      "student_id": 2,
      "status": "late",
      "remark": "มาสาย 10 นาที"
    },
    {
      "student_id": 3,
      "status": "absent",
      "remark": "ป่วย"
    }
  ]
}
```

**สถานะที่รองรับ:**
- `present` - มาเรียน
- `absent` - ขาดเรียน
- `late` - มาสาย
- `leave` - ลา

**Response:**
```json
{
  "success": true,
  "message": "บันทึกการเช็คชื่อสำเร็จ",
  "data": {
    "session_date": "2025-10-23",
    "classroom_name": "ม.1/1",
    "total_recorded": 30,
    "present": 27,
    "absent": 2,
    "late": 1,
    "leave": 0
  }
}
```

### 9. ประวัติการเช็คชื่อ

**Endpoint:** `GET /api/v1/attendance/history`  
**Authentication:** Required

**Query Parameters:**
- `classroom_id` (optional): รหัสห้องเรียน
- `date_from` (optional): วันที่เริ่มต้น (YYYY-MM-DD)
- `date_to` (optional): วันที่สิ้นสุด (YYYY-MM-DD)
- `page` (optional): หน้าที่ต้องการ
- `limit` (optional): จำนวนต่อหน้า

**Response:**
```json
{
  "success": true,
  "message": "ดึงข้อมูลสำเร็จ",
  "data": [
    {
      "session_date": "2025-10-23",
      "classroom": {
        "id": 1,
        "name": "ม.1/1",
        "grade": "ม.1"
      },
      "total_students": 30,
      "present": 28,
      "absent": 1,
      "late": 1,
      "leave": 0
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

---

## 📖 วิธีการใช้งาน

### 1. Authentication Flow

```javascript
// 1. เข้าสู่ระบบ
const loginResponse = await fetch('/api/v1/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'teacher@example.com',
    password: 'password123'
  })
});

const loginData = await loginResponse.json();

if (loginData.success) {
  // เก็บ token สำหรับใช้ในคำขอต่อไป
  localStorage.setItem('auth_token', loginData.token);
  localStorage.setItem('teacher_data', JSON.stringify(loginData.data));
}
```

### 2. การใช้ Token

```javascript
// ใช้ token ในทุกคำขอที่ต้องการ authentication
const token = localStorage.getItem('auth_token');

const response = await fetch('/api/v1/dashboard', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
```

### 3. Error Handling

```javascript
const response = await fetch('/api/v1/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(loginData)
});

const result = await response.json();

if (!result.success) {
  // แสดง error message
  alert(result.message);
  console.error(result.error);
} else {
  // ดำเนินการต่อ
  console.log('Success:', result.data);
}
```

---

## 💻 ตัวอย่างโค้ด

### React/JavaScript Frontend

```javascript
// utils/api.js
const API_BASE_URL = 'http://localhost:8080/api/v1/frontend';

class ApiService {
  constructor() {
    this.token = localStorage.getItem('auth_token');
  }

  async request(endpoint, options = {}) {
    const url = `${API_BASE_URL}${endpoint}`;
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    if (this.token) {
      config.headers['Authorization'] = `Bearer ${this.token}`;
    }

    try {
      const response = await fetch(url, config);
      const data = await response.json();

      if (!data.success) {
        throw new Error(data.message || 'เกิดข้อผิดพลาด');
      }

      return data;
    } catch (error) {
      console.error('API Error:', error);
      throw error;
    }
  }

  // Authentication
  async login(email, password) {
    const response = await this.request('/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });

    if (response.success) {
      this.token = response.token;
      localStorage.setItem('auth_token', response.token);
      localStorage.setItem('teacher_data', JSON.stringify(response.data));
    }

    return response;
  }

  async register(userData) {
    return this.request('/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  }

  // Dashboard
  async getDashboard() {
    return this.request('/dashboard');
  }

  // Classroom
  async createClassroom(name, grade) {
    return this.request('/classrooms', {
      method: 'POST',
      body: JSON.stringify({ name, grade }),
    });
  }

  async getClassroomDetail(id) {
    return this.request(`/classrooms/${id}`);
  }

  // Students
  async addStudent(studentData) {
    return this.request('/students', {
      method: 'POST',
      body: JSON.stringify(studentData),
    });
  }

  async getClassroomStudents(classroomId, page = 1, limit = 50) {
    return this.request(`/classrooms/${classroomId}/students?page=${page}&limit=${limit}`);
  }

  // Attendance
  async takeAttendance(attendanceData) {
    return this.request('/attendance', {
      method: 'POST',
      body: JSON.stringify(attendanceData),
    });
  }

  async getAttendanceHistory(params = {}) {
    const queryString = new URLSearchParams(params).toString();
    return this.request(`/attendance/history?${queryString}`);
  }
}

export default new ApiService();
```

### React Component ตัวอย่าง

```jsx
// components/Dashboard.jsx
import React, { useState, useEffect } from 'react';
import ApiService from '../utils/api';

function Dashboard() {
  const [dashboardData, setDashboardData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadDashboard();
  }, []);

  const loadDashboard = async () => {
    try {
      setLoading(true);
      const response = await ApiService.getDashboard();
      setDashboardData(response.data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div>กำลังโหลด...</div>;
  if (error) return <div>เกิดข้อผิดพลาด: {error}</div>;
  if (!dashboardData) return <div>ไม่พบข้อมูล</div>;

  return (
    <div className="dashboard">
      <h1>หน้าแรก</h1>
      
      {/* สรุปข้อมูล */}
      <div className="summary">
        <div className="card">
          <h3>ห้องเรียนทั้งหมด</h3>
          <p>{dashboardData.summary.total_classrooms}</p>
        </div>
        <div className="card">
          <h3>นักเรียนทั้งหมด</h3>
          <p>{dashboardData.summary.total_students}</p>
        </div>
        <div className="card">
          <h3>มาเรียนวันนี้</h3>
          <p>{dashboardData.summary.present_today}</p>
        </div>
        <div className="card">
          <h3>ขาดเรียนวันนี้</h3>
          <p>{dashboardData.summary.absent_today}</p>
        </div>
      </div>

      {/* รายการห้องเรียน */}
      <div className="classrooms">
        <h2>ห้องเรียนของฉัน</h2>
        {dashboardData.classrooms.map(classroom => (
          <div key={classroom.id} className="classroom-card">
            <h3>{classroom.name} ({classroom.grade})</h3>
            <p>นักเรียน: {classroom.student_count} คน</p>
            <button onClick={() => handleTakeAttendance(classroom.id)}>
              เช็คชื่อ
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}

export default Dashboard;
```

### ตัวอย่างการเช็คชื่อ

```jsx
// components/TakeAttendance.jsx
import React, { useState, useEffect } from 'react';
import ApiService from '../utils/api';

function TakeAttendance({ classroomId }) {
  const [students, setStudents] = useState([]);
  const [attendanceData, setAttendanceData] = useState({});
  const [sessionDate] = useState(new Date().toISOString().split('T')[0]);

  useEffect(() => {
    loadStudents();
  }, [classroomId]);

  const loadStudents = async () => {
    try {
      const response = await ApiService.getClassroomStudents(classroomId);
      setStudents(response.data);
      
      // ตั้งค่าเริ่มต้นเป็น present ทั้งหมด
      const initialData = {};
      response.data.forEach(student => {
        initialData[student.id] = {
          student_id: student.id,
          status: 'present',
          remark: ''
        };
      });
      setAttendanceData(initialData);
    } catch (error) {
      console.error('Error loading students:', error);
    }
  };

  const handleStatusChange = (studentId, status) => {
    setAttendanceData(prev => ({
      ...prev,
      [studentId]: {
        ...prev[studentId],
        status
      }
    }));
  };

  const handleRemarkChange = (studentId, remark) => {
    setAttendanceData(prev => ({
      ...prev,
      [studentId]: {
        ...prev[studentId],
        remark
      }
    }));
  };

  const handleSubmit = async () => {
    try {
      const attendanceArray = Object.values(attendanceData);
      
      const response = await ApiService.takeAttendance({
        classroom_id: classroomId,
        session_date: sessionDate,
        students: attendanceArray
      });

      alert(response.message);
      console.log('Attendance saved:', response.data);
    } catch (error) {
      alert('เกิดข้อผิดพลาด: ' + error.message);
    }
  };

  return (
    <div className="take-attendance">
      <h2>เช็คชื่อวันที่ {sessionDate}</h2>
      
      <div className="student-list">
        {students.map(student => (
          <div key={student.id} className="student-row">
            <div className="student-info">
              <span>{student.student_no}</span>
              <span>{student.first_name} {student.last_name}</span>
            </div>
            
            <div className="attendance-controls">
              <label>
                <input
                  type="radio"
                  name={`status_${student.id}`}
                  value="present"
                  checked={attendanceData[student.id]?.status === 'present'}
                  onChange={() => handleStatusChange(student.id, 'present')}
                />
                มาเรียน
              </label>
              
              <label>
                <input
                  type="radio"
                  name={`status_${student.id}`}
                  value="absent"
                  checked={attendanceData[student.id]?.status === 'absent'}
                  onChange={() => handleStatusChange(student.id, 'absent')}
                />
                ขาดเรียน
              </label>
              
              <label>
                <input
                  type="radio"
                  name={`status_${student.id}`}
                  value="late"
                  checked={attendanceData[student.id]?.status === 'late'}
                  onChange={() => handleStatusChange(student.id, 'late')}
                />
                มาสาย
              </label>
              
              <label>
                <input
                  type="radio"
                  name={`status_${student.id}`}
                  value="leave"
                  checked={attendanceData[student.id]?.status === 'leave'}
                  onChange={() => handleStatusChange(student.id, 'leave')}
                />
                ลา
              </label>
            </div>
            
            <input
              type="text"
              placeholder="หมายเหตุ"
              value={attendanceData[student.id]?.remark || ''}
              onChange={(e) => handleRemarkChange(student.id, e.target.value)}
              className="remark-input"
            />
          </div>
        ))}
      </div>
      
      <button onClick={handleSubmit} className="submit-btn">
        บันทึกการเช็คชื่อ
      </button>
    </div>
  );
}

export default TakeAttendance;
```

---

## 🚨 ข้อมูลสำคัญ

### Rate Limiting
- **Login/Register:** 10 requests/minute
- **API ทั่วไป:** 100 requests/minute

### Authentication
- Token มีอายุ 24 ชั่วโมง
- ต้องใส่ `Authorization: Bearer {token}` ในทุกคำขอที่ต้องการ authentication

### Error Codes
- `400` - Bad Request (ข้อมูลไม่ถูกต้อง)
- `401` - Unauthorized (ไม่มีสิทธิ์)
- `404` - Not Found (ไม่พบข้อมูล)
- `429` - Too Many Requests (เกิน rate limit)
- `500` - Internal Server Error (ข้อผิดพลาดของเซิร์ฟเวอร์)

### Data Validation
- Email ต้องเป็นรูปแบบอีเมลที่ถูกต้อง
- Password ต้องมีอย่างน้อย 6 ตัวอักษร
- ชื่อ-นามสกุล ต้องมีอย่างน้อย 2 ตัวอักษร
- วันที่ต้องเป็นรูปแบบ YYYY-MM-DD

---

**จัดทำโดย:** GitHub Copilot  
**ติดต่อ:** [API Documentation Team]  
**อัพเดทครั้งล่าสุด:** 23 ตุลาคม 2025