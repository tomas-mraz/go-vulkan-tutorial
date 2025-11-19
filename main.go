package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/tomas-mraz/vulkan"
)

var (
	window         *glfw.Window
	instance       vk.Instance
	physicalDevice vk.PhysicalDevice
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
		panic(err)
	}

	vk.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())
	err = vk.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err = glfw.CreateWindow(WIDTH, HEIGHT, "VulkanApp", nil, nil)
	if err != nil {
		panic("failed to create window")
	}
	fmt.Println("initWindow OK")
}

func initVulkan() {
	createInstance()
	setupDebugMessenger()
	pickPhysicalDevice()
	fmt.Println("initVulkan OK")
}

func pickPhysicalDevice() {
	var deviceCount uint32
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
	var deviceProperties vk.PhysicalDeviceProperties
	vk.GetPhysicalDeviceProperties(physicalDevice, &deviceProperties)
	deviceProperties.Deref()
	fmt.Println("Using device:", GetCString(deviceProperties.DeviceName[:]))
}

func isDeviceSuitable(device vk.PhysicalDevice) bool {
	var deviceProperties vk.PhysicalDeviceProperties
	var deviceFeatures vk.PhysicalDeviceFeatures
	vk.GetPhysicalDeviceProperties(device, &deviceProperties)
	vk.GetPhysicalDeviceFeatures(device, &deviceFeatures)
	deviceProperties.Deref()
	deviceFeatures.Deref()

	return deviceProperties.DeviceType == vk.PhysicalDeviceTypeIntegratedGpu && deviceFeatures.GeometryShader != 0
}

func setupDebugMessenger() {
	//TODO validation layers
}

func createInstance() {
	appInfo := vk.ApplicationInfo{
		SType:              vk.StructureTypeApplicationInfo,
		PApplicationName:   "Hello triangle",
		ApplicationVersion: vk.MakeVersion(1, 0, 0),
		PEngineName:        "No Engine",
		EngineVersion:      vk.MakeVersion(1, 0, 0),
		ApiVersion:         vk.ApiVersion10,
	}

	extensions := window.GetRequiredInstanceExtensions()
	createInfo := vk.InstanceCreateInfo{
		SType:                   vk.StructureTypeInstanceCreateInfo,
		PApplicationInfo:        &appInfo,
		EnabledExtensionCount:   uint32(len(extensions)),
		PpEnabledExtensionNames: extensions,
		EnabledLayerCount:       0,
	}

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
