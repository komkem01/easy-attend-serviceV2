# 🚀 Easy Attend Service - Quick Import Guide

## การ Import Postman Collection แบบเร็ว

### 1. Copy URL นี้และ Import ใน Postman:
```
https://raw.githubusercontent.com/komkem01/easy-attend-serviceV2/main/postman/Easy-Attend-Service-Complete-Collection.json
```

### 2. หรือ Import ไฟล์ที่อยู่ในเครื่อง:
- `Easy-Attend-Service-Complete-Collection.json`
- `Easy-Attend-Environment.json`

## 🎯 การทดสอบแบบเร็ว

### เริ่มต้น Server:
```bash
cd easy-attend-serviceV2
air
```

### ทดสอบ Complete Flow:
1. Import Collection + Environment
2. เลือก Environment "Easy Attend - Local Development"
3. Run folder "🚀 Complete Test Flow"
4. ตรวจสอบผลการทดสอบ

## ✅ สิ่งที่ควรได้หลังการทดสอบ:
- ✅ Teacher ลงทะเบียนสำเร็จ
- ✅ Login ได้ JWT token
- ✅ สร้าง Student, Classroom, Attendance สำเร็จ
- ✅ Logs บันทึกทุก action อัตโนมัติ
- ✅ Logout สำเร็จ

Happy Testing! 🎉