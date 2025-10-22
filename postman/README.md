# 🚀 Easy Attend Service - Postman Testing Guide

## � **อัปเดตล่าสุด - Complete System Migration (23 Oct 2025)**

✅ **การแก้ไขครบวงจร - ระบบพร้อมใช้งาน 100%:**
- **Database Migration**: เปลี่ยน Primary Keys จาก UUID → Integer IDs ทุก tables ✅
- **Model Consistency**: Teacher และ Student ใช้ `firstname`/`lastname` ทุกที่ ✅
- **Request Structure**: แก้ไข JSON field names ให้สอดคล้องกัน ✅
- **Controller Updates**: แก้ไข parameter handling สำหรับ Integer IDs ✅
- **Postman Collection**: อัปเดต requests ให้ใช้ `firstname`/`lastname` ✅
- **Attendance Status**: ใช้ภาษาอังกฤษ (`present`, `absent`, `late`, `leave`) ✅

🎯 **ทดสอบแล้ว**: Database, Models, Controllers, Requests, และ Postman Collection ทำงานสอดคล้องกันครบถ้วน!

## 📂 ไฟล์ที่จำเป็น
- `Easy-Attend-Service-Complete-Collection.json` - Postman Collection หลัก
- `Easy-Attend-Environment.json` - Environment variables

## 🔧 การติดตั้งใน Postman

### 1. Import Collection
1. เปิด Postman
2. คลิก **Import** 
3. เลือกไฟล์ `Easy-Attend-Service-Complete-Collection.json`
4. คลิก **Import**

### 2. Import Environment
1. คลิก **Import** อีกครั้ง
2. เลือกไฟล์ `Easy-Attend-Environment.json`
3. คลิก **Import**
4. เลือก Environment "Easy Attend - Local Development" ที่มุมขวาบน

## 🎯 วิธีการทดสอบ

### 🚀 Quick Start - Complete Test Flow
สำหรับการทดสอบแบบครบวงจรในครั้งเดียว:

1. **เริ่มต้น Server**
   ```bash
   air
   ```

2. **Run Complete Test Flow Folder**
   - ไปที่ folder "🚀 Complete Test Flow"
   - คลิกขวาที่ folder → **Run folder**
   - ระบบจะทำการทดสอบทั้งหมดตามลำดับ:
     - ลงทะเบียนครู (ชื่อสุ่มอัตโนมัติ)
     - สร้างนักเรียน (ชื่อสุ่มอัตโนมัติ)
     - สร้างห้องเรียน (ชื่อสุ่มอัตโนมัติ)
     - บันทึกการเข้าเรียน
     - ตรวจสอบ logs ทั้งหมด

🎲 **ข้อมูลสุ่มอัตโนมัติ**
- ชื่อครู: สุ่มจากชื่อไทย + email + เบอร์โทร
- ชื่อโรงเรียน: สุ่มจากโรงเรียนจริงในไทย  
- ชื่อนักเรียน: สุ่มจากชื่อไทย + รหัสนักเรียน
- ชื่อห้องเรียน: สุ่มรูปแบบ "ม.X/Y" หรือ "ป.X/Y"

### 📋 Step by Step Testing

⚠️ **สำคัญ: ต้องทำตามลำดับนี้เสมอ!**

#### 1. 🔐 Authentication (บังคับ - ต้องทำก่อน!)
```  
📁 🏫 Authentication
├── Register Teacher      (POST /auth/register) ← เริ่มตรงนี้เสมอ!
├── Login Teacher         (POST /auth/login) 
├── Get Profile          (GET /auth/profile)
└── Logout               (POST /auth/logout)
```

**เริ่มต้น:**
1. **Register Teacher** - สร้างบัญชีครูใหม่
2. **Login Teacher** - เข้าสู่ระบบ (จะได้ token อัตโนมัติ)

#### 2. 👥 จัดการข้อมูลผู้ใช้
```
📁 👨‍🏫 Teachers
├── Get All Teachers
├── Create Teacher
├── Get Teacher by ID
└── Update Teacher

📁 👩‍🎓 Students  
├── Get All Students
├── Create Student
├── Create Student 2
├── Get Student by ID
└── Update Student
```

#### 3. 🎯 ข้อมูลพื้นฐาน
```
📁 🎯 Lookup Data
├── Get All Genders
├── Get All Prefixes
├── Create Gender
└── Create Prefix
```

#### 4. 🏛️ จัดการห้องเรียน
```
📁 🏛️ Classrooms
├── Get All Classrooms
├── Create Classroom
├── Get Classroom by ID
└── Update Classroom

📁 👥 Classroom Members
├── Get All Classroom Members
├── Get Members by Classroom
├── Add Student to Classroom
└── Add Teacher to Classroom
```

#### 5. 📋 บันทึกการเข้าเรียน
```
📁 📋 Attendances
├── Get All Attendances
├── Create Attendance - Present
├── Create Attendance - Absent
├── Get Attendances by Classroom
└── Get Attendances by Student
```

#### 6. 📊 ตรวจสอบ Logs
```
📁 📊 Logs (Read-Only)
├── Get All Logs
├── Get Logs by Teacher
├── Get Logs by Action - Login
├── Get Logs by Action - Attendance
└── Create Manual Log (if needed)
```

## 🔄 Automatic Variables
Collection จะจัดการ variables อัตโนมัติ:

- `{{auth_token}}` - JWT token (จาก login)
- `{{teacher_id}}` - Teacher ID (จาก register/login)
- `{{student_id}}` - Student ID (จาก create student)
- `{{school_id}}` - School ID (จาก register)
- `{{classroom_id}}` - Classroom ID (จาก create classroom)
- `{{gender_id}}` - Gender ID (จาก get genders)
- `{{prefix_id}}` - Prefix ID (จาก get prefixes)

## 📊 ข้อมูลทดสอบตัวอย่าง

### ครู (Teachers)
```json
{
    "email": "teacher1@school.com",
    "password": "password123",
    "firstname": "สมชาย",
    "lastname": "ใจดี",
    "phone": "081-234-5678",
    "school_name": "โรงเรียนวัดไผ่"
}
```

### นักเรียน (Students)
```json
{
    "student_no": "STD001",
    "firstname": "นางสาวสมใส",
    "lastname": "เรียนดี",
    "school_name": "โรงเรียนวัดไผ่"
}
```

### ห้องเรียน (Classrooms)
```json
{
    "school_id": 1,
    "teacher_id": 1,
    "name": "ม.3/1",
    "grade": "ม.3"
}
```

### การเข้าเรียน (Attendance)
```json
{
    "classroom_id": 1,
    "teacher_id": 1,
    "student_id": 1,
    "session_date": "2025-10-23",
    "status": "present",
    "checked_at": 1729656000,
    "remark": "เข้าเรียนตรงเวลา"
}
```

## 🔍 สิ่งที่ควรตรวจสอบ

### ✅ Automatic Logging
หลังจากทำ action ต่างๆ ให้ตรวจสอบ logs:

1. **Login** → ตรวจสอบ `/logs/action?action=login`
2. **Create Student** → ตรวจสอบ `/logs/action?action=create_student`
3. **Create Classroom** → ตรวจสอบ `/logs/action?action=create_classroom`
4. **Attendance** → ตรวจสอบ `/logs/action?action=attendance`
5. **Logout** → ตรวจสอบ `/logs/action?action=logout`

### ✅ Response Status Codes
- **200** - Success (GET, PUT)
- **201** - Created (POST)
- **400** - Bad Request
- **401** - Unauthorized
- **404** - Not Found
- **500** - Internal Server Error

### ✅ Authentication Flow
1. Register → ได้ teacher data
2. Login → ได้ JWT token
3. ใช้ token ในการเรียก protected endpoints
4. Logout → ลบ session และ log activity

## 🐛 การแก้ไขปัญหา

### ❌ "Invalid UUID length: 0" 
**สาเหตุ**: ไม่ได้ทำ Authentication ก่อน
**วิธีแก้**: 
1. Run **Register Teacher** ก่อนเสมอ (ระบบจะสุ่มชื่อใหม่อัตโนมัติ)
2. จากนั้น Run **Login Teacher** 
3. ตรวจสอบว่า Environment variables ได้ค่าแล้ว (teacher_id, school_id)

### ❌ "Name already exists" 
**สาเหตุ**: ข้อมูลซ้ำกับที่สร้างไว้แล้ว
**วิธีแก้**: ระบบจะสุ่มชื่อใหม่อัตโนมัติทุกครั้ง ไม่ต้องกังวล!

### ❌ Token หมดอายุ
- Run **Login Teacher** ใหม่เพื่อรับ token ใหม่

### ❌ Server ไม่ตอบสนอง
- ตรวจสอบว่า server รันอยู่ที่ port 8080
- ตรวจสอบ database connection

### ❌ Variables ไม่ได้ค่า
- ตรวจสอบว่าเลือก Environment ที่ถูกต้อง
- **ต้อง Run Register Teacher ก่อนเสมอ!**
- Run requests ตามลำดับเพื่อให้ได้ variables

## 📈 การทดสอบ Performance
สามารถใช้ **Collection Runner** เพื่อ:
- Run ทั้ง collection พร้อมกัน
- ทดสอบ load หลาย iterations
- Export ผลการทดสอบ

## 🎉 Happy Testing!
Collection นี้ครอบคลุมการทดสอบทุกฟีเจอร์ของ Easy Attend Service รวมถึงระบบ Automatic Logging ที่จะบันทึกทุก action อัตโนมัติ!