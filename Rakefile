class VersionedFile
  def initialize(file, regex)
    @file = file
    @regex = regex
  end

  def current_version!
    @current_version ||= matched_data![1]
  end

  def bump_version!(type)
    position = case type
               when :major
                 0
               when :minor
                 1
               when :patch
                 2
               end
    @current_version = current_version!.split('.').tap do |v|
      v[position] = v[position].to_i + 1
      # Reset consequent numbers
      ((position + 1)..2).each { |p| v[p] = 0 }
    end.join('.')
  end

  def save!
    text = File.read(@file)
    new_line = matched_data![0].gsub(matched_data![1], @current_version)
    text.gsub!(matched_data![0], new_line)

    File.open(@file, 'w') { |f| f.puts text }
  end

  private

  def matched_data!
    @matched_data ||= begin
                        m = @regex.match File.read(@file)
                        raise "No version #{@regex} matched in #{@file}" unless m
                        m
                      end
  end
end

def fullpath(file)
  File.expand_path(file, File.dirname(__FILE__))
end

VERSION_FILES = {
  fullpath('commands/version.go') => /^const Version = "(\d+.\d+.\d+)"$/,
  fullpath('README.md')           => /Current version is \[(\d+.\d+.\d+)\]/,
  fullpath('.goxc.json')          => /"PackageVersion": "(\d+.\d+.\d+)"/,
  fullpath('homebrew/gh.rb')      => /VERSION = '(\d+.\d+.\d+)'/
}

namespace :release do
  desc "Current released version"
  task :current do
    vf = VersionedFile.new(*VERSION_FILES.first)
    puts vf.current_version!
  end

  [:major, :minor, :patch].each do |type|
    desc "Release #{type} version"
    task type do
      VERSION_FILES.each do |file, regex|
        begin
          vf = VersionedFile.new(file, regex)
          current_version = vf.current_version!
          vf.bump_version!(type)
          vf.save!
          puts "Successfully bump #{file} from #{current_version} to #{vf.current_version!}"
        rescue => e
          puts e
        end
      end
    end
  end
end
