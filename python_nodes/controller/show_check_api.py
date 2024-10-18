import dlib
import cv2
import os
from flask import Blueprint, request, jsonify

# Đường dẫn tới file shape_predictor_68_face_landmarks.dat
current_dir = os.path.dirname(os.path.abspath(__file__))
PREDICTOR_PATH = os.path.join(current_dir, "shape_predictor_68_face_landmarks.dat")

# Khởi tạo dlib's face detector (HOG-based) và facial landmarks predictor
detector = dlib.get_frontal_face_detector()
predictor = dlib.shape_predictor(PREDICTOR_PATH)



show_check_bp = Blueprint('show_check', __name__)

def draw_face_box(image_path, output_path):
    # Đọc ảnh
    img = cv2.imread(image_path)
    
    # Chuyển đổi sang ảnh grayscale
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    
    # Phát hiện khuôn mặt trong ảnh
    faces = detector(gray)
    
    for face in faces:
        # Dự đoán các điểm landmark trên khuôn mặt
        landmarks = predictor(gray, face)
        
        # Tọa độ của vùng lông mày (phần trên) và cằm (phần dưới)
        eyebrow_points = []
        chin_points = []
        
        # Lông mày trái (từ điểm 17 đến 21)
        for i in range(17, 22):
            x = landmarks.part(i).x
            y = landmarks.part(i).y
            eyebrow_points.append((x, y))
        
        # Lông mày phải (từ điểm 22 đến 26)
        for i in range(22, 27):
            x = landmarks.part(i).x
            y = landmarks.part(i).y
            eyebrow_points.append((x, y))
        
        # Cằm (từ điểm 6 đến 11)
        for i in range(6, 12):
            x = landmarks.part(i).x
            y = landmarks.part(i).y
            chin_points.append((x, y))
        
        # Tìm điểm trái và phải xa nhất từ các điểm lông mày và cằm
        left = min(eyebrow_points + chin_points, key=lambda p: p[0])[0]
        right = max(eyebrow_points + chin_points, key=lambda p: p[0])[0]
        top = min(eyebrow_points, key=lambda p: p[1])[1]
        bottom = max(chin_points, key=lambda p: p[1])[1]
        
        # Vẽ box từ vùng lông mày xuống cằm
        cv2.rectangle(img, (left, top), (right, bottom), (0, 255, 0), 2)
    
    # Lưu ảnh đã chỉnh sửa ra file mới
    cv2.imwrite(output_path, img)

@show_check_bp.route('/show_check', methods=['POST'])
def process_image():
    # Nhận dữ liệu JSON từ request
    data = request.json
    
    if 'input_image_path' not in data:
        print(1)
        return jsonify({'error': 'No input_image_path provided'}), 400
    if 'save_path' not in data:
        print(2)
        return jsonify({'error': 'No save_path provided'}), 400
    
    input_image_path = data['input_image_path']
    save_path = data["save_path"]
    
    if not os.path.exists(input_image_path):
        return jsonify({'error': 'Input image path does not exist'}), 400
    
    # Xử lý ảnh: Vẽ box từ lông mày xuống cằm
    draw_face_box(input_image_path, save_path)
    
    
    # Trả về file đã xử lý
    return jsonify({ "result": True })
