# Copyright 2020 Erik Maciejewski
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.load("@bazel_tools//tools/build_defs/cc:action_names.bzl", "ACTION_NAMES")

load(
    "@bazel_tools//tools/cpp:cc_toolchain_config_lib.bzl",
    "action_config",
    "feature",
    "feature_set",
    "flag_group",
    "flag_set",
    "make_variable",
    "tool",
    "tool_path",
    "with_feature_set",
)

load(
    "@bazel_tools//tools/cpp:cc_toolchain_config.bzl",
    ALL_COMPILE_ACTIONS = "all_compile_actions",
    ALL_CPP_COMPILE_ACTIONS = "all_cpp_compile_actions",
    ALL_LINK_ACTIONS = "all_link_actions",
)

def _linux_arm64_impl(ctx):
    toolchain_identifier = "gcc-linux-arm64-cross"
    compiler = "gcc"
    abi_version = "elf"
    abi_libc_version = "glibc_unknown"
    target_libc = "glibc_unknown"
    target_cpu = ctx.attr.target.split("-")[0]
    root = "/usr/xcc/"
    arch = "aarch64-unknown-linux-gnueabi"
    install_path = root + arch + "/"
    include_path_prefix = install_path + arch + "/"
    sysroot = include_path_prefix + "/sysroot/"

    # TODO: do we need to explicitly call out systemroot includes?
    cross_system_include_dirs = [
        include_path_prefix + "include/c++/4.9.4",
        include_path_prefix + "include/c++/4.9.4/" + arch,
        install_path + "lib/gcc/aarch64-unknown-linux-gnueabi/4.9.4/include",
    ]

    # TODO: do we need to explicitly call out systemroot libs?
    cross_system_lib_dirs = [
        include_path_prefix + "lib",
        install_path + "lib",
    ]

    opt_feature = feature(name = "opt")
    dbg_feature = feature(name = "dbg")
    fastbuild_feature = feature(name = "fastbuild")
    random_seed_feature = feature(name = "random_seed", enabled = True)
    supports_pic_feature = feature(name = "supports_pic", enabled = True)
    supports_dynamic_linker_feature = feature(name = "supports_dynamic_linker", enabled = True)

    # TODO: add the applicable features
    features = []

    tool_paths = [
        tool_path(name = "ld", path = install_path +"bin/" + arch + "-ld"),
        tool_path(name = "cpp", path = install_path +"bin/" + arch + "-cpp"),
        tool_path(name = "dwp", path = install_path +"bin/" + arch + "-dwp"),
        tool_path(name = "gcov", path = install_path +"bin/" + arch + "-gcov"),
        tool_path(name = "nm", path = install_path +"bin/" + arch + "-nm"),
        tool_path(name = "objcopy", path = install_path +"bin/" + arch + "-objcopy"),
        tool_path(name = "objdump", path = install_path +"bin/" + arch + "-objdump"),
        tool_path(name = "strip", path = install_path +"bin/" + arch + "-strip"),
        tool_path(name = "gcc", path = install_path +"bin/" + arch + "-gcc"),
        tool_path(name = "ar", path = install_path +"bin/" + arch + "-ar"),
    ]

    return cc_common.create_cc_toolchain_config_info(
        ctx = ctx,
        features = features,
        abi_version = abi_version,
        builtin_sysroot = sysroot,
        compiler = compiler,
        cxx_builtin_include_directories = cross_system_include_dirs,
        host_system_name = "x86_64-unknown-linux-gnu",
        target_cpu = target_cpu,
        target_libc = target_libc,
        target_system_name = ctx.attr.target,
        tool_paths = tool_paths,
        toolchain_identifier = toolchain_identifier,
    )

# mvp for arm only, windows and osx later
cc_toolchain_config = rule(
    implementation = _linux_arm64_impl,
    attrs = {
        "target": attr.string(mandatory = True),
        "stdlib": attr.string(),
    },
    provides = [CcToolchainConfigInfo],
)
