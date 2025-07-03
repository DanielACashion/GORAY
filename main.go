package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Screen struct {
	Width  int32
	Height int32
}

const MAX_COLUMNS int32 = 20
const CAMERA_MOUSE_MOVE_SENSITIVITY float32 = 0.003
const CAMERA_ROTATION_SPEED float32 = 0.03
const CAMERA_MOVE_SPEED float32 = 5.4
const CAMERA_PAN_SPEED float32 = 0.2
const CAMERA_ORBITAL_SPEED float32 = 0.5

func main() {
	var screen Screen = Screen{Width: 800, Height: 450}
	var frames int64 = 0
	var zoomLevel int32 = 0
	_ = zoomLevel

	//var frametime float32 = 0

	var camera rl.Camera
	camera.Position = rl.Vector3{X: 0.0, Y: 2.0, Z: 4.0}
	camera.Target = rl.Vector3{X: 0.0, Y: 2.0, Z: 0.0}
	camera.Up = rl.Vector3{X: 0.0, Y: 1.0, Z: 0.0}
	camera.Fovy = 60
	cameraMode := rl.CameraFirstPerson

	//fontBm := rl.LoadFont("pixantiqua.fnt")
	//msg := "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHI\nJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmn\nopqrstuvwxyz" +
	//"{|}~¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓ\nÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷\nøùúûüýþÿ"

	rl.InitWindow(screen.Width, screen.Height, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.DisableCursor()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		fmt.Println(rl.GetMouseDelta())
		//update
		UpdateCamera(&camera, cameraMode)
		//frametime = rl.GetFrameTime()
		frames += 1

		//RENDER
		rl.BeginDrawing()
		rl.BeginMode3D(camera)
		rl.ClearBackground(rl.White)

		rl.DrawPlane(rl.Vector3{X: 0.0, Y: 2.0, Z: 4.0}, rl.Vector2{X: 32.0, Y: 32.0}, rl.LightGray) // Draw ground
		rl.DrawCube(rl.Vector3{X: -16.0, Y: 2.5, Z: 0.0}, 1.0, 5.0, 32.0, rl.Blue)                   // Draw a blue wall
		rl.DrawCube(rl.Vector3{X: 16.0, Y: 2.5, Z: 0.0}, 1.0, 5.0, 32.0, rl.Lime)                    // Draw a green wall
		rl.DrawCube(rl.Vector3{X: 0.0, Y: 2.5, Z: 16.0}, 32.0, 5.0, 1.0, rl.Yellow)                  // Draw a yellow wall
		rl.EndMode3D()

		rl.DrawRectangle(5, 5, 330, 100, rl.Fade(rl.SkyBlue, 0.5))
		rl.DrawRectangleLines(5, 5, 330, 100, rl.Blue)
		rl.DrawText("Camera controls:", 15, 15, 10, rl.Black)
		rl.DrawText("- Move keys: W, A, S, D, Space, Left-Ctrl", 15, 30, 10, rl.Black)
		rl.DrawText("- Look around: arrow keys or mouse", 15, 45, 10, rl.Black)
		rl.DrawText("- Camera mode keys: 1, 2, 3, 4", 15, 60, 10, rl.Black)
		rl.DrawText("- Zoom keys: num-plus, num-minus or mouse scroll", 15, 75, 10, rl.Black)
		rl.DrawText("- Camera projection key: P", 15, 90, 10, rl.Black)

		rl.DrawRectangle(600, 5, 195, 100, rl.Fade(rl.SkyBlue, 0.5))
		rl.DrawRectangleLines(600, 5, 195, 100, rl.Blue)
		rl.DrawText("Camera status:", 610, 15, 10, rl.Black)
		// rl.DrawText(TextFormat("- Mode: %s", (cameraMode == CAMERA_FREE) ? "FREE" :
		//                                       (cameraMode == CAMERA_FIRST_PERSON) ? "FIRST_PERSON" :
		//                                       (cameraMode == CAMERA_THIRD_PERSON) ? "THIRD_PERSON" :
		//                                       (cameraMode == CAMERA_ORBITAL) ? "ORBITAL" : "CUSTOM"), 610, 30, 10, BLACK);
		// rl.DrawText(TextFormat("- Projection: %s", (camera.projection == CAMERA_PERSPECTIVE) ? "PERSPECTIVE" :
		//                                         (camera.projection == CAMERA_ORTHOGRAPHIC) ? "ORTHOGRAPHIC" : "CUSTOM"), 610, 45, 10, BLACK);
		rl.DrawText(fmt.Sprintf("- Position: (%06.3f, %06.3f, %06.3f)", camera.Position.X, camera.Position.Y, camera.Position.Z), 610, 60, 10, rl.Black)
		rl.DrawText(fmt.Sprintf("- Target: (%06.3f, %06.3f, %06.3f)", camera.Target.X, camera.Target.Y, camera.Target.Z), 610, 75, 10, rl.Black)
		rl.DrawText(fmt.Sprintf("- Up: (%06.3f, %06.3f, %06.3f)", camera.Up.X, camera.Up.Y, camera.Up.Z), 610, 90, 10, rl.Black)

		// rl.DrawTextEx(fontBm, fmt.Sprintf("%.2f", 1/frametime), rl.Vector2{X: 20.0, Y: 0}, float32(fontBm.BaseSize), 2, rl.Gray)
		// rl.DrawTextEx(fontBm, msg, rl.Vector2{X: 20.0, Y: float32(rl.GetScreenHeight()) - 30}, float32(fontBm.BaseSize), 2, rl.Gray)
		//rl.DrawText("Using BMFont (Angelcode) imported", 20, int32(rl.GetScreenHeight())-30, 20, rl.Gray)
		fmt.Println(frames)
		rl.EndDrawing()
	}
}

// local version of the rl cam, accounts for false gamepad thanks to keychron
func UpdateCamera(camera *rl.Camera, mode rl.CameraMode) {
	var mousePositionDelta = rl.GetMouseDelta()

	//vars for later
	moveInWorldPlaneBool := mode == rl.CameraFirstPerson || mode == rl.CameraThirdPerson
	var moveInWorldPlane uint8
	if moveInWorldPlaneBool {
		moveInWorldPlane = 1
	}

	rotateAroundTargetBool := mode == rl.CameraThirdPerson || mode == rl.CameraOrbital
	var rotateAroundTarget uint8
	if rotateAroundTargetBool {
		rotateAroundTarget = 1
	}

	lockViewBool := mode == rl.CameraFirstPerson || mode == rl.CameraThirdPerson || mode == rl.CameraOrbital
	var lockView uint8
	if lockViewBool {
		lockView = 1
	}
	var rotateUp uint8

	var cameraOrbitSpeed float32 = CAMERA_ORBITAL_SPEED * rl.GetFrameTime()
	var cameraRotateSpeed float32 = CAMERA_ROTATION_SPEED * rl.GetFrameTime()
	var cameraPanSpeed float32 = CAMERA_PAN_SPEED * rl.GetFrameTime()
	var cameraMoveSpeed float32 = CAMERA_MOVE_SPEED * rl.GetFrameTime()

	if mode == rl.CameraCustom {
		return
	} else if mode == rl.CameraOrbital {
		var rotation rl.Matrix = rl.MatrixRotate(rl.GetCameraUp(camera), cameraOrbitSpeed)
		var view rl.Vector3 = rl.Vector3Subtract(camera.Position, camera.Target)
		view = rl.Vector3Transform(view, rotation)
		camera.Position = rl.Vector3Add(camera.Position, view)
	} else {
		if rl.IsKeyDown(rl.KeyDown) {
			rl.CameraPitch(camera, -cameraRotateSpeed, lockView, rotateAroundTarget, rotateUp)
		}
		if rl.IsKeyDown(rl.KeyUp) {
			rl.CameraPitch(camera, cameraRotateSpeed, lockView, rotateAroundTarget, rotateUp)
		}
		if rl.IsKeyDown(rl.KeyRight) {
			rl.CameraYaw(camera, -cameraRotateSpeed, rotateAroundTarget)
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			rl.CameraYaw(camera, cameraRotateSpeed, rotateAroundTarget)
		}
		if rl.IsKeyDown(rl.KeyQ) {
			rl.CameraRoll(camera, -cameraRotateSpeed)
		}
		if rl.IsKeyDown(rl.KeyE) {
			rl.CameraRoll(camera, cameraRotateSpeed)
		}

		// Camera movement
		if mode == rl.CameraFree && rl.IsMouseButtonDown(rl.MouseButtonMiddle) {

			if mousePositionDelta.X > 0.0 {
				rl.CameraMoveRight(camera, cameraPanSpeed, moveInWorldPlane)
			}
			if mousePositionDelta.X < 0.0 {
				rl.CameraMoveRight(camera, -cameraPanSpeed, moveInWorldPlane)
			}
			if mousePositionDelta.Y > 0.0 {
				rl.CameraMoveUp(camera, cameraPanSpeed)
			}
			if mousePositionDelta.Y < 0.0 {
				rl.CameraMoveUp(camera, -cameraPanSpeed)
			}
		} else {
			// Mouse support
			rl.CameraYaw(camera, -mousePositionDelta.X*CAMERA_MOUSE_MOVE_SENSITIVITY, rotateAroundTarget)
			rl.CameraPitch(camera, -mousePositionDelta.Y*CAMERA_MOUSE_MOVE_SENSITIVITY, lockView, rotateAroundTarget, rotateUp)
		}
		if rl.IsKeyDown(rl.KeyW) {
			rl.CameraMoveForward(camera, cameraMoveSpeed, moveInWorldPlane)
		}
		if rl.IsKeyDown(rl.KeyS) {
			rl.CameraMoveForward(camera, -cameraMoveSpeed, moveInWorldPlane)
		}
		if rl.IsKeyDown(rl.KeyD) {
			rl.CameraMoveRight(camera, cameraMoveSpeed, moveInWorldPlane)
		}
		if rl.IsKeyDown(rl.KeyA) {
			rl.CameraMoveRight(camera, -cameraMoveSpeed, moveInWorldPlane)
		}
		if rl.IsGamepadAvailable(0) && rl.GetGamepadName(0) != "Keychron Link " {
			gpname, isAvail := rl.GetGamepadName(0), rl.IsGamepadAvailable(0)
			fmt.Println(gpname)
			_ = isAvail
			rl.CameraYaw(camera, -(rl.GetGamepadAxisMovement(0, rl.GamepadAxisRightX)*2)*CAMERA_MOUSE_MOVE_SENSITIVITY, rotateAroundTarget)
			rl.CameraPitch(camera, -(rl.GetGamepadAxisMovement(0, rl.GamepadAxisRightY)*2)*CAMERA_MOUSE_MOVE_SENSITIVITY, lockView, rotateAroundTarget, rotateUp)

			if rl.GetGamepadAxisMovement(0, rl.GamepadAxisLeftY) >= -0.25 {
				rl.CameraMoveForward(camera, cameraMoveSpeed, moveInWorldPlane)
			}
			if rl.GetGamepadAxisMovement(0, rl.GamepadAxisLeftY) >= 0.25 {
				rl.CameraMoveForward(camera, -cameraMoveSpeed, moveInWorldPlane)
			}
			if rl.GetGamepadAxisMovement(0, rl.GamepadAxisLeftX) >= -0.25 {
				rl.CameraMoveRight(camera, cameraMoveSpeed, moveInWorldPlane)
			}
			if rl.GetGamepadAxisMovement(0, rl.GamepadAxisLeftX) >= 0.25 {
				rl.CameraMoveRight(camera, -cameraMoveSpeed, moveInWorldPlane)
			}
		}
		if mode == rl.CameraFree {
			if rl.IsKeyDown(rl.KeySpace) {
				rl.CameraMoveUp(camera, cameraMoveSpeed)
			}
			if rl.IsKeyDown(rl.KeyLeftControl) {
				rl.CameraMoveUp(camera, -cameraMoveSpeed)
			}
		}
	}
	if mode == rl.CameraThirdPerson || mode == rl.CameraOrbital || mode == rl.CameraFree {
		rl.CameraMoveToTarget(camera, -rl.GetMouseWheelMove())
		if rl.IsKeyPressed(rl.KeyKpSubtract) {
			rl.CameraMoveToTarget(camera, 2.0)
		}
		if rl.IsKeyPressed(rl.KeyKpAdd) {
			rl.CameraMoveToTarget(camera, -2.0)
		}
	}

}
