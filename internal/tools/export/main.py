import os
import tkinter as tk
from tkinter import filedialog, messagebox

from ExportUtil import ExportUtil


def select_directory():
    directory = filedialog.askdirectory()
    if directory:
        directory_label.config(text=directory)


def export():
    directory = directory_label.cget("text")
    if directory == "未选择目录":
        messagebox.showwarning("警告", "请先选择目录")
    else:
        # 在这里添加你的导出逻辑
        ExportUtil.exportData(directory)
        messagebox.showinfo("成功", f"已导出到目录: {directory}")


def createDirectory(directory):
    # 检查目录是否存在，如果不存在则创建
    if not os.path.exists(directory):
        try:
            os.makedirs(directory)  # 创建目录
        except Exception as e:
            pass


createDirectory("json")
createDirectory("ts")
# 创建主窗口
root = tk.Tk()
root.title("导表工具")

# 固定窗口宽高
window_width = 400  # 窗口宽度
window_height = 200  # 窗口高度
root.geometry(f"{window_width}x{window_height}")  # 设置窗口大小
root.resizable(False, False)  # 禁止调整窗口大小

# 居中显示窗口
screen_width = root.winfo_screenwidth()  # 获取屏幕宽度
screen_height = root.winfo_screenheight()  # 获取屏幕高度
x = (screen_width - window_width) // 2  # 计算窗口左上角的 x 坐标
y = (screen_height - window_height) // 2  # 计算窗口左上角的 y 坐标
root.geometry(f"+{x}+{y}")  # 设置窗口位置

# 创建选择目录的按钮
select_button = tk.Button(root, text="选择excel目录", command=select_directory)
select_button.pack(pady=10)

# 显示选择的目录
directory_label = tk.Label(root, text="未选择目录", fg="blue")
directory_label.pack(pady=10)

# 创建导出按钮
export_button = tk.Button(root, text="导出", command=export)
export_button.pack(pady=10)

# 运行主循环
root.mainloop()
