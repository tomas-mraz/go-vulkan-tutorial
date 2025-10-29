package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/tomas-mraz/vulkan"
)

var (
	window         *glfw.Window
	instance       *vk.Instance
	physicalDevice *vk.PhysicalDevice
)

const (
	WIDTH  = 800
	HEIGHT = 600
)

func init() {
	runtime.LockOSThread()
}

func main() {
	initWindow()
	initVulkan()
	mainLoop()
	cleanup()
}

func initWindow() {
	var err error
	err = glfw.Init()
	if err != nil {
		panic("failed GLFW initialization")
	}
	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err = glfw.CreateWindow(WIDTH, HEIGHT, "VulkanApp", nil, nil)
	if err != nil {
		panic("failed to create window")
	}
}

func initVulkan() {
	createInstance()
	setupDebugMessenger()
	pickPhysicalDevice()
}

func pickPhysicalDevice() {
	deviceCount := uint32(0)
	result := vk.EnumeratePhysicalDevices(instance, &deviceCount, nil)
	if result != vk.Success {
		panic("failed to enumerate physical devices")
	}
	fmt.Println("Found", deviceCount, "devices")
	if deviceCount == 0 {
		panic("failed to find GPUs with Vulkan support")
	}

	devices := make([]vk.PhysicalDevice, deviceCount)
	result = vk.EnumeratePhysicalDevices(instance, &deviceCount, devices)
	if result != vk.Success {
		panic("failed to enumerate physical devices")
	}
	for _, device := range devices {
		if isDeviceSuitable(device) {
			physicalDevice = device
			break
		}
	}
	if physicalDevice == nil {
		panic("failed to find a suitable GPU")
	}
	fmt.Println("Using device:", physicalDevice.Properties.DeviceName)
}

func isDeviceSuitable(device vk.PhysicalDevice) bool {
	//return true
	var deviceProperties vk.PhysicalDeviceProperties
	var deviceFeatures vk.PhysicalDeviceFeatures
	vk.GetPhysicalDeviceProperties(device, &deviceProperties)
	vk.GetPhysicalDeviceFeatures(device, &deviceFeatures)

	return deviceProperties.DeviceType == vk.PhysicalDeviceTypeIntegratedGpu && deviceFeatures.GeometryShader != 0
}

func setupDebugMessenger() {
	//TODO validation layers
}

func createInstance() {
	var appInfo vk.ApplicationInfo
	appInfo.SType = vk.StructureTypeApplicationInfo
	appInfo.PApplicationName = "Hello triangle"
	appInfo.ApplicationVersion = vk.MakeVersion(1, 0, 0)
	appInfo.PEngineName = "No Engine"
	appInfo.EngineVersion = vk.MakeVersion(1, 0, 0)
	appInfo.ApiVersion = vk.ApiVersion10

	var createInfo vk.InstanceCreateInfo
	createInfo.SType = vk.StructureTypeInstanceCreateInfo
	createInfo.PApplicationInfo = &appInfo

	extensions := window.GetRequiredInstanceExtensions()
	createInfo.EnabledExtensionCount = uint32(len(extensions))
	createInfo.PpEnabledExtensionNames = extensions
	createInfo.EnabledLayerCount = 0

	result := vk.CreateInstance(&createInfo, nil, &instance)
	if result != vk.Success {
		panic("failed to create Vulkan instance")
	}
	fmt.Println("Instance created")
}

func mainLoop() {
}

func cleanup() {
	vk.DestroyInstance(instance, nil)
	window.Destroy()
	glfw.Terminate()
}
