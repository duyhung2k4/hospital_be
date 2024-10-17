from flask import Blueprint, request, jsonify
import face_recognition
import os
import numpy as np
from PIL import ImageEnhance, Image

face_encoding_bp = Blueprint('face_encoding', __name__)

def adjust_brightness(image, factor):
    """Điều chỉnh độ sáng của ảnh."""
    pil_image = Image.fromarray(image)
    enhancer = ImageEnhance.Brightness(pil_image)
    enhanced_image = enhancer.enhance(factor)  # Tăng hoặc giảm độ sáng
    return np.array(enhanced_image)

@face_encoding_bp.route('/face_encoding', methods=['POST'])
def face_encoding():
    directory_path = request.json.get("directory_path")
    
    if not directory_path or not os.path.exists(directory_path):
        return jsonify({"result": "error", "message": "Directory path is invalid."}), 400

    list_face_encoding = []
    errors = []

    # Các mức độ sáng khác nhau: 5 mức, 1.0 là sáng trung bình (ảnh gốc)
    brightness_factors = [0.6, 0.8, 1.0, 1.2, 1.6]

    try:
        for image_file in os.listdir(directory_path):
            image_path = os.path.join(directory_path, image_file)
            
            try:
                new_image = face_recognition.load_image_file(image_path)
                
                # Tăng cường dữ liệu với các mức độ sáng khác nhau
                for factor in brightness_factors:
                    # Điều chỉnh độ sáng của ảnh
                    brightened_image = adjust_brightness(new_image, factor)

                    # Lấy vị trí khuôn mặt
                    face_locations = face_recognition.face_locations(brightened_image)
                    
                    # Lấy các điểm đặc trưng trên khuôn mặt
                    face_landmarks_list = face_recognition.face_landmarks(brightened_image)

                    if len(face_locations) == 0 or len(face_landmarks_list) == 0:
                        errors.append(f"No face found in {image_file} with brightness factor {factor}")
                        continue

                    # Lặp qua từng khuôn mặt trong hình ảnh
                    for i, face_location in enumerate(face_locations):
                        # Lấy các điểm lông mày
                        landmarks = face_landmarks_list[i]
                        left_eyebrow = landmarks.get("left_eyebrow", [])
                        right_eyebrow = landmarks.get("right_eyebrow", [])

                        if left_eyebrow and right_eyebrow:
                            # Tính tọa độ `y` trung bình của lông mày để xác định khu vực từ lông mày trở xuống
                            eyebrow_top_y = min([point[1] for point in left_eyebrow + right_eyebrow])

                            # Điều chỉnh vị trí khuôn mặt để chỉ lấy từ trên lông mày trở xuống
                            top, right, bottom, left = face_location
                            new_top = eyebrow_top_y  # Tọa độ mới từ trên lông mày
                            
                            # Cập nhật vùng khuôn mặt để mã hóa chỉ từ lông mày trở xuống
                            new_face_location = (new_top, right, bottom, left)

                            # Mã hóa đặc trưng khuôn mặt
                            new_face_encodings = face_recognition.face_encodings(brightened_image, [new_face_location])

                            if len(new_face_encodings) > 0:
                                new_face_encoding = new_face_encodings[0]
                                list_face_encoding.append(new_face_encoding.tolist())
                            else:
                                errors.append(f"No face encoding found in {image_file} with brightness factor {factor}")

            except Exception as e:
                errors.append(f"Error processing {image_file}: {str(e)}")

        return jsonify({"result": "success", "face_encodings": list_face_encoding, "errors": errors})

    except Exception as e:
        return jsonify({"result": "error", "message": str(e)}), 500
